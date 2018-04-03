package mongodb

import (
	"errors"

	"gopkg.in/mgo.v2"
)

type QueryParams struct {
	Id         interface{}
	Query      interface{}
	Selector   interface{}
	Skip       int
	Limit      int
	SortFields []string
}

var ErrNotQueryParams error = errors.New("parameter is NOT QueryParams")

func NewQueryParams() *QueryParams {
	return new(QueryParams)
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
	return q
}

func ConvertToQueryParams(params interface{}) (
	queryParams *QueryParams, err error) {
	if params == nil {
		return nil, nil
	}
	pqp, ok := params.(*QueryParams)
	if ok {
		return pqp, nil
	}
	qp, ok := params.(QueryParams)
	if ok {
		return &qp, nil
	}
	return nil, ErrNotQueryParams
}
