package mongodb

import (
	"github.com/donyori/cqa/common/json"
)

type MongoDbSettings struct {
	Url    string            `json:"url"`
	DbName string            `json:"db_name"`
	CNames map[string]string `json:"collection_names"`
}

const (
	MgoCNameKeyQ   string = "q"
	MgoCNameKeyTag string = "tag"
	MgoCNameKeyQlf string = "qlf"
	MgoCNameKeyQv  string = "qv"

	MgoSettingsFilename string = "settings/mgo.json"
)

var (
	GlobalSettings MongoDbSettings
)

func init() {
	// Default values:
	GlobalSettings.Url = "127.0.0.1:27017"
	GlobalSettings.DbName = "cqa"
	GlobalSettings.CNames = map[string]string{
		MgoCNameKeyQ:   "questions.v1",
		MgoCNameKeyTag: "tags.v1",
		MgoCNameKeyQlf: "question_linguistic_features.v1",
		MgoCNameKeyQv:  "question_vector.v1",
	}

	_, err := json.DecodeJsonFromFileIfExist(MgoSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
