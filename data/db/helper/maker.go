package helper

import (
	"errors"

	"github.com/donyori/cqa/data/db/id"
	mhelper "github.com/donyori/cqa/data/model/helper"
)

var ErrNoCorrespondingMaker error = errors.New(
	"cannot find corresponding maker")

func GetMakerByCollectionId(cid id.CollectionId) (
	maker mhelper.Maker, err error) {
	if !cid.IsValid() {
		return nil, id.ErrInvalidCollectionId
	}
	maker = nil
	err = nil
	switch cid {
	case id.QuestionCollection:
		maker = mhelper.QuestionMaker
	case id.QuestionVectorCollection:
		maker = mhelper.QuestionVectorMaker
	default:
		err = ErrNoCorrespondingMaker
	}
	return maker, err
}
