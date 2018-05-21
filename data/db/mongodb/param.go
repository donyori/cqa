package mongodb

import (
	"errors"
	"reflect"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db/generic"
)

type QueryParams struct {
	Id         interface{}
	Query      interface{}
	Selector   interface{}
	Skip       int
	Limit      int
	SortFields []string
	Hint       []string
}

var ErrNotQueryParams error = errors.New("parameter is NOT QueryParams")

func NewParamById(id interface{}) (param bson.M, err error) {
	if id == nil {
		return nil, generic.ErrNilId
	}
	return bson.M{"_id": id}, nil
}

func NewParamByIds(ids interface{}) (param bson.M, err error) {
	if ids == nil {
		return nil, generic.ErrEmptyIds
	}
	v := reflect.ValueOf(ids)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil, generic.ErrIdsNotSlice
	}
	if v.Kind() != reflect.Slice {
		return nil, generic.ErrIdsNotSlice
	}
	if v.Len() == 0 {
		return nil, generic.ErrEmptyIds
	}
	return bson.M{"_id": bson.M{"$in": v.Interface()}}, nil
}

func NewQueryParams() *QueryParams {
	return new(QueryParams)
}

func NewQueryParamsById(id interface{}) (qp *QueryParams, err error) {
	if id == nil {
		return nil, generic.ErrNilId
	}
	qp = NewQueryParams()
	qp.Id = id
	return qp, nil
}

func NewQueryParamsByIds(ids interface{}) (qp *QueryParams, err error) {
	if ids == nil {
		return nil, generic.ErrEmptyIds
	}
	v := reflect.ValueOf(ids)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil, generic.ErrIdsNotSlice
	}
	if v.Kind() != reflect.Slice {
		return nil, generic.ErrIdsNotSlice
	}
	if v.Len() == 0 {
		return nil, generic.ErrEmptyIds
	}
	qp = NewQueryParams()
	qp.Query = bson.M{"_id": bson.M{"$in": v.Interface()}}
	qp.Limit = v.Len()
	return qp, nil
}

func (params *QueryParams) MakeQuery(c *mgo.Collection) *mgo.Query {
	if c == nil {
		return nil
	}
	if params == nil {
		return c.Find(nil)
	}
	var q *mgo.Query
	if params.Id != nil {
		q = c.FindId(params.Id)
	} else {
		q = c.Find(params.Query)
	}
	if params.Selector != nil {
		q = q.Select(params.Selector)
	}
	if params.Skip > 0 {
		q = q.Skip(params.Skip)
	}
	if params.Limit > 0 {
		q = q.Limit(params.Limit)
	}
	if len(params.SortFields) > 0 {
		q = q.Sort(params.SortFields...)
	}
	if len(params.Hint) > 0 {
		q = q.Hint(params.Hint...)
	}
	return q
}

func ConvertToQueryParams(params interface{}) (
	queryParams *QueryParams, err error) {
	if params == nil {
		return nil, nil
	}
	v := reflect.ValueOf(params)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil, ErrNotQueryParams
	}
	if v.CanAddr() {
		v = v.Addr()
		params = v.Interface()
		qp, ok := params.(*QueryParams)
		if ok {
			return qp, nil
		}
	} else {
		params = v.Interface()
		qpv, ok := params.(QueryParams)
		if ok {
			return &qpv, nil
		}
	}
	return nil, ErrNotQueryParams
}
