package helper

import (
	"errors"
	"reflect"

	"github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model"
)

var ErrNoCorrespondingModelType error = errors.New(
	"cannot find corresponding model type")

func GetModelTypeByCollectionId(cid id.CollectionId) (
	modelType reflect.Type, err error) {
	if !cid.IsValid() {
		return nil, id.ErrInvalidCollectionId
	}
	modelType = nil
	err = nil
	switch cid {
	case id.QuestionCollection:
		modelType = reflect.TypeOf(model.Question{})
	case id.QuestionVectorCollection:
		modelType = reflect.TypeOf(model.QuestionVector{})
	default:
		err = ErrNoCorrespondingModelType
	}
	return modelType, err
}
