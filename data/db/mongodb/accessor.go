package mongodb

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db/generic"
	dbhelper "github.com/donyori/cqa/data/db/helper"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model/helper"
)

type MgoAccessor struct {
	WithMgoSession
}

var ErrNilAccessor error = errors.New("MgoAccessor is nil")

func NewMgoAccessor(session generic.Session) (
	accessor *MgoAccessor, err error) {
	accessor = new(MgoAccessor)
	if session != nil {
		err = accessor.SetSession(session)
		if err != nil {
			return nil, err
		}
	}
	return accessor, nil
}

func (ma *MgoAccessor) Get(cid dbid.CollectionId, params interface{},
	maker helper.Maker) (res interface{}, err error) {
	if ma == nil {
		return nil, ErrNilAccessor
	}
	if maker == nil {
		maker, err = dbhelper.GetMakerByCollectionId(cid)
		if err != nil {
			return nil, err
		}
	}
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return nil, err
	}
	if qp == nil {
		qp = NewQueryParams()
	}
	session, c, err := ma.aquireSessionAndCollection(cid)
	if err != nil {
		return nil, err
	}
	defer session.Release()
	qp.Limit = 1
	q := qp.MakeQuery(c)
	res = maker()
	err = q.One(res)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (ma *MgoAccessor) GetById(cid dbid.CollectionId, id interface{},
	maker helper.Maker) (res interface{}, err error) {
	if ma == nil {
		return nil, ErrNilAccessor
	}
	if id == nil {
		return nil, generic.ErrNilId
	}
	params := NewQueryParams()
	params.Id = id
	return ma.Get(cid, params, maker)
}

func (ma *MgoAccessor) Scan(cid dbid.CollectionId, params interface{},
	bufferSize int, quitC <-chan struct{}, maker helper.Maker) (
	outC <-chan interface{}, resC <-chan error, err error) {
	if ma == nil {
		return nil, nil, ErrNilAccessor
	}
	if maker == nil {
		maker, err = dbhelper.GetMakerByCollectionId(cid)
		if err != nil {
			return nil, nil, err
		}
	}
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return nil, nil, err
	}
	session, c, err := ma.aquireSessionAndCollection(cid)
	if err != nil {
		return nil, nil, err
	}
	// Session release in the scan goroutine.
	var out chan interface{}
	if bufferSize > 0 {
		out = make(chan interface{}, bufferSize)
	} else {
		out = make(chan interface{})
	}
	res := make(chan error, 1)
	go func() {
		defer session.Release()
		q := qp.MakeQuery(c)
		iter := q.Iter()
		defer iter.Close()
		defer close(res)
		defer close(out)
		result := maker()
		isQuit := false
		for !isQuit && iter.Next(result) {
			select {
			case <-quitC:
				isQuit = true
			default:
				out <- result
				result = maker() // Make new one each time.
			}
		}
		res <- iter.Err()
	}()
	return out, res, nil
}

func (ma *MgoAccessor) Save(cid dbid.CollectionId, selector interface{},
	model interface{}) (isNew bool, err error) {
	if ma == nil {
		return false, ErrNilAccessor
	}
	session, c, err := ma.aquireSessionAndCollection(cid)
	if err != nil {
		return false, err
	}
	defer session.Release()
	if selector == nil {
		id, _ := helper.GetId(model)
		if id != nil {
			selector = bson.M{"_id": id}
		} else {
			err = c.Insert(model)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}
	info, err := c.Upsert(selector, model)
	if err != nil {
		return false, err
	}
	return info.Updated == 0, nil
}

func (ma *MgoAccessor) SaveById(cid dbid.CollectionId, id interface{},
	model interface{}) (isNew bool, err error) {
	if ma == nil {
		return false, ErrNilAccessor
	}
	if id == nil {
		return false, generic.ErrNilId
	}
	return ma.Save(cid, bson.M{"_id": id}, model)
}
