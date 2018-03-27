package db

import (
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/model"
)

type QuestionAccessor interface {
	Connector

	Get(params interface{}) (question *model.Question, err error)
	GetById(id interface{}) (question *model.Question, err error)
	Scan(params interface{}, bufferSize int) (out <-chan *model.Question, res <-chan error, quit chan<- struct{}, err error)
	Save(question *model.Question) (isNew bool, err error)
}

type QuestionVectorAccessor interface {
	Connector

	Get(params interface{}) (questionVector *model.QuestionVector, err error)
	GetById(id interface{}) (questionVector *model.QuestionVector, err error)
	Scan(params interface{}, bufferSize int) (out <-chan *model.QuestionVector, res <-chan error, quit chan<- struct{}, err error)
	Save(questionVector *model.QuestionVector) (isNew bool, err error)
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
