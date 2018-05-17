package id

import "errors"

type CollectionId int8

const (
	MetaCollection CollectionId = iota
	QuestionCollection
	QuestionVectorCollection
	QuestionClassificationCollection
)

var (
	ValidCollectionIds = [...]CollectionId{
		MetaCollection,
		QuestionCollection,
		QuestionVectorCollection,
		QuestionClassificationCollection,
	}

	ErrInvalidCollectionId error = errors.New("collection ID is invalid")

	collectionIdStrings = [...]string{
		"Meta",
		"Question",
		"Question Vector",
		"Question Classification",
	}
)

func (cid CollectionId) IsValid() bool {
	return cid >= MetaCollection && cid <= QuestionClassificationCollection
}

func (cid CollectionId) String() string {
	if !cid.IsValid() {
		return "Unknown"
	}
	return collectionIdStrings[cid]
}
