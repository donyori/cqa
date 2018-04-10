package helper

import (
	"errors"

	"github.com/donyori/cqa/data/model/tag"
)

var (
	ErrNilModel         error = errors.New("model is nil")
	ErrUnknownModelType error = errors.New("model type is unknown")
	ErrNoSuchField      error = errors.New("such field does NOT exist")
)

func GetId(model interface{}) (id interface{}, err error) {
	return getIdBase(model, tag.DataModelTagKey, tag.DataModelTagId)
}

func GetMongoDbId(model interface{}) (id interface{}, err error) {
	return getIdBase(model, tag.MongoDbTagKey, tag.MongoDbTagId)
}

func getIdBase(model interface{}, tagKey string, tagValue string) (
	id interface{}, err error) {
	if model == nil {
		return nil, ErrNilModel
	}
	id, ok := tag.GetDataModelFieldValueByTag(model, tagKey, tagValue)
	if !ok {
		return nil, ErrNoSuchField
	}
	return id, nil
}
