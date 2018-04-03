package mongodb

import (
	"errors"
	"sync"

	"gopkg.in/mgo.v2"

	"github.com/donyori/cqa/data/db/generic"
	"github.com/donyori/cqa/data/db/id"
)

type MgoSession struct {
	sess     *mgo.Session
	settings *MongoDbSettings
	// "url" field is important to ensure the URL cannot be changed
	// after NewNewMgoSession().
	url string

	acquireCounter int
	cond           *sync.Cond
}

type WithMgoSession struct {
	sess *MgoSession
	lock sync.RWMutex
}

type sessionSeed struct {
	sess      *mgo.Session
	copyCount int
}

var (
	ErrPoolLimitExceeded    error = errors.New("pool limit exceeded")
	ErrInvalidSession       error = errors.New("sess is invalid")
	ErrSessionPoolCleanedUp error = errors.New("session pool is cleaned up")
	ErrNilWithSession       error = errors.New("WithMgoSession is nil")
	ErrNotMgoSession        error = errors.New("session is not a MgoSession")

	sessionSeedMap map[string]*sessionSeed
	mapLock        sync.Mutex
)

func init() {
	sessionSeedMap = make(map[string]*sessionSeed)
}

func acquireSession(url string, settings *MongoDbSettings) (
	sess *mgo.Session, err error) {
	if settings == nil {
		settings = &GlobalSettings
	}
	mapLock.Lock()
	defer mapLock.Unlock()
	if sessionSeedMap == nil {
		return nil, ErrSessionPoolCleanedUp
	}
	seed := sessionSeedMap[url]
	if seed == nil {
		var newSess *mgo.Session
		newSess, err = mgo.Dial(url)
		if err != nil || newSess == nil {
			return nil, err
		}
		seed = &sessionSeed{
			sess:      newSess,
			copyCount: 0,
		}
		sessionSeedMap[url] = seed
	}

	if settings.PoolLimit > 0 {
		// Check limit first
		if seed.copyCount >= settings.PoolLimit {
			return nil, ErrPoolLimitExceeded
		}
	}

	sess = seed.sess.Copy()
	err = sess.Ping()
	if err != nil {
		sess.Close()
		return nil, err
	}
	seed.copyCount++
	return sess, nil
}

func releaseSession(url string, sess *mgo.Session) {
	if sess == nil {
		return
	}
	mapLock.Lock()
	defer mapLock.Unlock()
	if sessionSeedMap == nil {
		return
	}
	seed := sessionSeedMap[url]
	if seed == nil {
		return
	}
	seed.copyCount--
}

func CleanUpSessionPool() {
	mapLock.Lock()
	defer mapLock.Unlock()
	if sessionSeedMap == nil {
		return
	}
	for _, seed := range sessionSeedMap {
		if seed != nil && seed.sess != nil {
			seed.sess.Close()
		}
	}
	sessionSeedMap = nil
}

func GetCollection(cid id.CollectionId, nativeSession interface{},
	settings interface{}) (c *mgo.Collection, err error) {
	if !cid.IsValid() {
		return nil, id.ErrInvalidCollectionId
	}
	sess, ok := nativeSession.(*mgo.Session)
	if !ok {
		return nil, ErrInvalidSession
	}
	mgoSettings, ok := settings.(*MongoDbSettings)
	if !ok {
		return nil, ErrInvalidSession
	}
	cName := mgoSettings.CNames[cid]
	if cName == "" {
		return nil, ErrCollectionNameNotSet
	}
	c = sess.DB(mgoSettings.DbName).C(cName)
	return c, nil
}

func NewMgoSession(settings *MongoDbSettings) (sess *MgoSession, err error) {
	if settings == nil {
		settings = &GlobalSettings
	}
	url := settings.Url
	nativeSess, err := acquireSession(url, settings)
	if err != nil {
		return nil, err
	}
	sess = &MgoSession{
		sess:     nativeSess,
		settings: settings,
		url:      url,
		cond:     sync.NewCond(new(sync.Mutex)),
	}
	return sess, nil
}

func (ms *MgoSession) IsValid() bool {
	if ms == nil || ms.cond == nil {
		return false
	}
	ms.cond.L.Lock()
	defer ms.cond.L.Unlock()
	return ms.isValid()
}

func (ms *MgoSession) GetUrl() (url string, err error) {
	if ms == nil || ms.cond == nil {
		return "", ErrInvalidSession
	}
	ms.cond.L.Lock()
	defer ms.cond.L.Unlock()
	if !ms.isValid() {
		return "", ErrInvalidSession
	}
	return ms.url, nil
}

func (ms *MgoSession) Acquire() (nativeSession interface{},
	settings interface{}, err error) {
	if ms == nil || ms.cond == nil {
		return nil, nil, ErrInvalidSession
	}
	ms.cond.L.Lock()
	defer ms.cond.L.Unlock()
	if !ms.isValid() {
		return nil, nil, ErrInvalidSession
	}
	ms.acquireCounter++
	return ms.sess, ms.settings, nil
}

func (ms *MgoSession) Release() {
	if ms == nil || ms.cond == nil {
		return
	}
	ms.cond.L.Lock()
	defer ms.cond.L.Unlock()
	if !ms.isValid() {
		return
	}
	ms.acquireCounter--
	ms.cond.Broadcast()
}

func (ms *MgoSession) Close() {
	if ms == nil || ms.cond == nil {
		return
	}
	ms.cond.L.Lock()
	defer ms.cond.L.Unlock()
	if ms.sess == nil {
		return
	}
	for ms.acquireCounter > 0 {
		ms.cond.Wait()
	}
	if ms.sess != nil {
		releaseSession(ms.url, ms.sess)
		ms.sess.Close()
		ms.sess = nil
	}
}

func (ms *MgoSession) isValid() bool {
	return ms.sess != nil && ms.sess.Ping() == nil
}

func (wms *WithMgoSession) GetSession() generic.Session {
	if wms == nil {
		return nil
	}
	wms.lock.RLock()
	defer wms.lock.RUnlock()
	return wms.sess
}

func (wms *WithMgoSession) SetSession(session generic.Session) error {
	if wms == nil {
		return ErrNilWithSession
	}
	sess, ok := session.(*MgoSession)
	if !ok {
		return ErrNotMgoSession
	}
	wms.lock.Lock()
	defer wms.lock.Unlock()
	wms.sess = sess
	return nil
}

func (wms *WithMgoSession) aquireSessionAndCollection(cid id.CollectionId) (
	session generic.Session, c *mgo.Collection, err error) {
	session = wms.GetSession()
	if session == nil {
		return nil, nil, ErrInvalidSession
	}
	nativeSession, settings, err := session.Acquire()
	if err != nil {
		return nil, nil, err
	}
	c, err = GetCollection(cid, nativeSession, settings)
	if err != nil {
		session.Release()
		return nil, nil, err
	}
	return session, c, nil
}
