package mongodb

import (
	"errors"
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db/generic"
	dbhelper "github.com/donyori/cqa/data/db/helper"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model/helper"
)

type Accessor struct {
	WithSession
}

var ErrNilAccessor error = errors.New("MongoDB accessor is nil")

func NewAccessor(session generic.Session) (
	accessor *Accessor, err error) {
	accessor = new(Accessor)
	if session != nil {
		err = accessor.SetSession(session)
		if err != nil {
			return nil, err
		}
	}
	return accessor, nil
}

func (a *Accessor) IsExisted(cid dbid.CollectionId, params interface{}) (
	res bool, err error) {
	if a == nil {
		return false, ErrNilAccessor
	}
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return false, err
	}
	if qp == nil {
		qp = NewQueryParams()
	}
	qp.Limit = 1
	n, err := a.Count(cid, qp)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (a *Accessor) IsExistedById(cid dbid.CollectionId, id interface{}) (
	res bool, err error) {
	if a == nil {
		return false, ErrNilAccessor
	}
	params, err := NewQueryParamsById(id)
	if err != nil {
		return false, err
	}
	return a.IsExisted(cid, params)
}

func (a *Accessor) FetchOne(cid dbid.CollectionId, params interface{},
	modelType reflect.Type) (res interface{}, err error) {
	if a == nil {
		return nil, ErrNilAccessor
	}
	if modelType == nil {
		modelType, err = dbhelper.GetModelTypeByCollectionId(cid)
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
	qp.Limit = 1
	session, c, err := a.aquireSessionAndCollection(cid)
	if err != nil {
		return nil, err
	}
	defer session.Release()
	q := qp.MakeQuery(c)
	res = reflect.New(modelType).Interface()
	err = q.One(res)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (a *Accessor) FetchOneById(cid dbid.CollectionId, id interface{},
	modelType reflect.Type) (res interface{}, err error) {
	if a == nil {
		return nil, ErrNilAccessor
	}
	params, err := NewQueryParamsById(id)
	if err != nil {
		return nil, err
	}
	return a.FetchOne(cid, params, modelType)
}

func (a *Accessor) FetchAll(cid dbid.CollectionId, params interface{},
	modelType reflect.Type) (res interface{}, err error) {
	if a == nil {
		return nil, ErrNilAccessor
	}
	if modelType == nil {
		modelType, err = dbhelper.GetModelTypeByCollectionId(cid)
		if err != nil {
			return nil, err
		}
	}
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return nil, err
	}
	session, c, err := a.aquireSessionAndCollection(cid)
	if err != nil {
		return nil, err
	}
	defer session.Release()
	q := qp.MakeQuery(c)
	spt := reflect.SliceOf(reflect.PtrTo(modelType))
	pResV := reflect.New(spt)
	err = q.All(pResV.Interface())
	if err != nil {
		return nil, err
	}
	return pResV.Elem().Interface(), nil
}

func (a *Accessor) FetchAllByIds(cid dbid.CollectionId, ids interface{},
	modelType reflect.Type) (res interface{}, err error) {
	if a == nil {
		return nil, ErrNilAccessor
	}
	params, err := NewQueryParamsByIds(ids)
	if err != nil {
		return nil, err
	}
	return a.FetchAll(cid, params, modelType)
}

func (a *Accessor) Scan(cid dbid.CollectionId, params interface{},
	bufferSize uint32, quitC <-chan struct{}, modelType reflect.Type) (
	outC <-chan interface{}, resC <-chan error, err error) {
	if a == nil {
		return nil, nil, ErrNilAccessor
	}
	if modelType == nil {
		modelType, err = dbhelper.GetModelTypeByCollectionId(cid)
		if err != nil {
			return nil, nil, err
		}
	}
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return nil, nil, err
	}
	session, c, err := a.aquireSessionAndCollection(cid)
	if err != nil {
		return nil, nil, err
	}
	// Session release in the scan goroutine.
	out := make(chan interface{}, bufferSize)
	res := make(chan error, 1)
	go func() {
		defer session.Release()
		q := qp.MakeQuery(c)
		iter := q.Iter()
		defer iter.Close()
		defer close(res)
		defer close(out)
		result := reflect.New(modelType).Interface()
		isQuit := false
		for !isQuit && iter.Next(result) {
			select {
			case <-quitC:
				isQuit = true
			default:
				out <- result
				// Make new one each time.
				result = reflect.New(modelType).Interface()
			}
		}
		res <- iter.Err()
	}()
	return out, res, nil
}

func (a *Accessor) ScanByIds(cid dbid.CollectionId, ids interface{},
	bufferSize uint32, quitC <-chan struct{}, modelType reflect.Type) (
	outC <-chan interface{}, resC <-chan error, err error) {
	if a == nil {
		return nil, nil, ErrNilAccessor
	}
	params, err := NewQueryParamsByIds(ids)
	if err != nil {
		return nil, nil, err
	}
	return a.Scan(cid, params, bufferSize, quitC, modelType)
}

func (a *Accessor) Count(cid dbid.CollectionId, params interface{}) (
	res int64, err error) {
	if a == nil {
		return 0, ErrNilAccessor
	}
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return 0, err
	}
	if qp == nil {
		qp = NewQueryParams()
	}
	qp.Selector = bson.M{"_id": 1}
	session, c, err := a.aquireSessionAndCollection(cid)
	if err != nil {
		return 0, err
	}
	defer session.Release()
	q := qp.MakeQuery(c)
	n, err := q.Count()
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

func (a *Accessor) SaveOne(cid dbid.CollectionId, selector interface{},
	model interface{}) (isNew bool, err error) {
	if a == nil {
		return false, ErrNilAccessor
	}
	session, c, err := a.aquireSessionAndCollection(cid)
	if err != nil {
		return false, err
	}
	defer session.Release()
	if selector == nil {
		id, _ := helper.GetMongoDbId(model)
		if id != nil {
			selector, err = NewParamById(id)
			if err != nil {
				return false, err
			}
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

func (a *Accessor) SaveOneById(cid dbid.CollectionId, id interface{},
	model interface{}) (isNew bool, err error) {
	if a == nil {
		return false, ErrNilAccessor
	}
	selector, err := NewParamById(id)
	if err != nil {
		return false, err
	}
	return a.SaveOne(cid, selector, model)
}
