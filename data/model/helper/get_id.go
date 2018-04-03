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
	if model == nil {
		return nil, ErrNilModel
	}
	id, ok := tag.GetDataModelFieldValueByTag(model, tag.DataModelTagId)
	if !ok {
		return nil, ErrNoSuchField
	}
	return id, nil
}
