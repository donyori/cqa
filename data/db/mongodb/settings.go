package mongodb

import (
	"errors"

	"github.com/donyori/cqa/common/json"
	"github.com/donyori/cqa/data/db/id"
)

type Settings struct {
	Url    string                     `json:"url"`
	DbName string                     `json:"db_name"`
	CNames map[id.CollectionId]string `json:"collection_names"`

	PoolLimit int `json:"pool_limit"`
}

const SettingsFilename string = "../settings/mongodb.json"

var (
	GlobalSettings Settings

	ErrCollectionNameNotSet error = errors.New("collection name is not set")
)

func init() {
	// Default values:
	GlobalSettings.Url = "127.0.0.1:27017"
	GlobalSettings.DbName = "cqa"
	GlobalSettings.CNames = map[id.CollectionId]string{
		id.MetaCollection:                   "meta.v1",
		id.QuestionCollection:               "question.v2",
		id.QuestionVectorCollection:         "question_vector.v2",
		id.QuestionClassificationCollection: "question_classification.v2",
	}
	GlobalSettings.PoolLimit = 1024

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
