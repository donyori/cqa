package db

import (
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/dtype"
)

type QuestionAccessor interface {
	Get(params interface{}) (question *dtype.Question, err error)
	GetById(id interface{}) (question *dtype.Question, err error)
	Scan(bufferSize int, params interface{}) (out <-chan *dtype.Question, res <-chan error, quit chan<- struct{}, err error)
}

func NewQuestionAccessor() (accessor QuestionAccessor, err error) {
	switch GlobalSettings.DbType {
	case DbTypeMongoDB:
		return new(mongodb.MgoQuestionAccessor), nil
	default:
		return nil, ErrUnknownDbType
	}
}
