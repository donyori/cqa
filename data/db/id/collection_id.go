package id

import "errors"

type CollectionId int8

const (
	MetaCollection CollectionId = iota
	QuestionCollection
	TagCollection
	QuestionLinguisticFeatureCollection
	QuestionVectorCollection
)

var (
	ValidCollectionIds = [...]CollectionId{
		MetaCollection,
		QuestionCollection,
		TagCollection,
		QuestionLinguisticFeatureCollection,
		QuestionVectorCollection,
	}

	ErrInvalidCollectionId error = errors.New("collection ID is invalid")

	collectionIdStrings = [...]string{
		"Meta",
		"Question",
		"Tag",
		"Question Linguistic Feature",
		"Question Vector",
	}
)

func (cid CollectionId) IsValid() bool {
	return cid >= MetaCollection && cid <= QuestionVectorCollection
}

func (cid CollectionId) String() string {
	if !cid.IsValid() {
		return "Unknown"
	}
	return collectionIdStrings[cid]
}
