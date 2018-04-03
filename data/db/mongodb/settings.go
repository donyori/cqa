package mongodb

import (
	"errors"

	"github.com/donyori/cqa/common/json"
	"github.com/donyori/cqa/data/db/id"
)

type MongoDbSettings struct {
	Url    string                     `json:"url"`
	DbName string                     `json:"db_name"`
	CNames map[id.CollectionId]string `json:"collection_names"`

	PoolLimit int `json:"pool_limit"`
}

const MgoSettingsFilename string = "settings/mgo.json"

var (
	GlobalSettings MongoDbSettings

	ErrCollectionNameNotSet error = errors.New("collection name is not set")
)

func init() {
	// Default values:
	GlobalSettings.Url = "127.0.0.1:27017"
	GlobalSettings.DbName = "cqa"
	GlobalSettings.CNames = map[id.CollectionId]string{
		id.QuestionCollection: "question.v1",
		id.TagCollection:      "tag.v1",
		id.QuestionLinguisticFeatureCollection: "" +
			"question_linguistic_feature.v1",
		id.QuestionVectorCollection: "question_vector.v1",
	}
	GlobalSettings.PoolLimit = 1024

	_, err := json.DecodeJsonFromFileIfExist(
		MgoSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
