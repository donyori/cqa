package db

import (
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/dtype"
)

type QuestionAccessor interface {
	Connector

	Get(params interface{}) (question *dtype.Question, err error)
	GetById(id interface{}) (question *dtype.Question, err error)
	Scan(params interface{}, bufferSize int) (out <-chan *dtype.Question, res <-chan error, quit chan<- struct{}, err error)
	Save(question *dtype.Question) (isNew bool, err error)
}

type QuestionVectorAccessor interface {
	Connector

	Get(params interface{}) (questionVector *dtype.QuestionVector, err error)
	GetById(id interface{}) (questionVector *dtype.QuestionVector, err error)
	Scan(params interface{}, bufferSize int) (out <-chan *dtype.QuestionVector, res <-chan error, quit chan<- struct{}, err error)
	Save(questionVector *dtype.QuestionVector) (isNew bool, err error)
}

func NewQuestionAccessor() (accessor QuestionAccessor, err error) {
	switch GlobalSettings.DbType {
	case DbTypeMongoDB:
		return mongodb.NewMgoQuestionAccessor(nil), nil
	default:
		return nil, ErrUnknownDbType
	}
}

func NewQuestionVectorAccessor() (accessor QuestionVectorAccessor, err error) {
	switch GlobalSettings.DbType {
	case DbTypeMongoDB:
		return mongodb.NewMgoQuestionVectorAccessor(nil), nil
	default:
		return nil, ErrUnknownDbType
	}
}
