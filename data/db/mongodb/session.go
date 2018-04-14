package mongodb

import (
	"errors"
	"sync"

	"gopkg.in/mgo.v2"

	"github.com/donyori/cqa/data/db/generic"
	"github.com/donyori/cqa/data/db/id"
)

type Session struct {
	sess     *mgo.Session
	settings *Settings
	// "url" field is important to ensure the URL cannot be changed
	// after NewSession().
	url string

	acquireCounter int
	cond           *sync.Cond
}

type WithSession struct {
	sess *Session
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
	ErrNilWithSession       error = errors.New("WithSession is nil")
	ErrNotSession           error = errors.New(
		"session is not a MongoDB session")

	sessionSeedMap     map[string]*sessionSeed
	sessionSeedMapLock sync.Mutex
)

func init() {
	sessionSeedMap = make(map[string]*sessionSeed)
}

func acquireSession(url string, settings *Settings) (
	sess *mgo.Session, err error) {
	if settings == nil {
		settings = &GlobalSettings
	}
	sessionSeedMapLock.Lock()
	defer sessionSeedMapLock.Unlock()
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
	sessionSeedMapLock.Lock()
	defer sessionSeedMapLock.Unlock()
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
	sessionSeedMapLock.Lock()
	defer sessionSeedMapLock.Unlock()
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
	mgoSettings, ok := settings.(*Settings)
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

func NewSession(settings *Settings) (sess *Session, err error) {
	if settings == nil {
		settings = &GlobalSettings
	}
	url := settings.Url
	nativeSess, err := acquireSession(url, settings)
	if err != nil {
		return nil, err
	}
	sess = &Session{
		sess:     nativeSess,
		settings: settings,
		url:      url,
		cond:     sync.NewCond(new(sync.Mutex)),
	}
	return sess, nil
}

func (ms *Session) IsValid() bool {
	if ms == nil || ms.cond == nil {
		return false
	}
	ms.cond.L.Lock()
	defer ms.cond.L.Unlock()
	return ms.isValid()
}

func (ms *Session) GetUrl() (url string, err error) {
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

func (ms *Session) Acquire() (nativeSession interface{},
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

func (ms *Session) Release() {
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

func (ms *Session) Close() {
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

func (ms *Session) isValid() bool {
	return ms.sess != nil && ms.sess.Ping() == nil
}

func (wms *WithSession) GetSession() generic.Session {
	if wms == nil {
		return nil
	}
	wms.lock.RLock()
	defer wms.lock.RUnlock()
	return wms.sess
}

func (wms *WithSession) SetSession(session generic.Session) error {
	if wms == nil {
		return ErrNilWithSession
	}
	sess, ok := session.(*Session)
	if !ok {
		return ErrNotSession
	}
	wms.lock.Lock()
	defer wms.lock.Unlock()
	wms.sess = sess
	return nil
}

func (wms *WithSession) aquireSessionAndCollection(cid id.CollectionId) (
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
