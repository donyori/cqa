package id

import "errors"

type CollectionId int8

const (
	QuestionCollection CollectionId = iota
	TagCollection
	QuestionLinguisticFeatureCollection
	QuestionVectorCollection
)

var (
	ValidCollectionIds = [...]CollectionId{
		QuestionCollection,
		TagCollection,
		QuestionLinguisticFeatureCollection,
		QuestionVectorCollection,
	}

	ErrInvalidCollectionId error = errors.New("collection ID is invalid")

	collectionIdStrings = [...]string{
		"Question",
		"Tag",
		"Question Linguistic Feature",
		"Question Vector",
	}
)

func (cid CollectionId) IsValid() bool {
	return cid >= QuestionCollection && cid <= QuestionVectorCollection
}

func (cid CollectionId) String() string {
	if !cid.IsValid() {
		return "Unknown"
	}
	return collectionIdStrings[cid]
}
