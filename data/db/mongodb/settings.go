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
	MgoCNameKeyQa  string = "qa"
	MgoCNameKeyTag string = "tag"
	MgoCNameKeyQlf string = "qlf"

	MgoSettingsFilename string = "settings/mgo.json"
)

var (
	GlobalSettings MongoDbSettings
)

func init() {
	// Default values:
	GlobalSettings.Url = "127.0.0.1"
	GlobalSettings.DbName = "cqa"
	GlobalSettings.CNames = map[string]string{
		MgoCNameKeyQa:  "qa",
		MgoCNameKeyTag: "tag",
		MgoCNameKeyQlf: "qlf",
	}

	_, err := json.DecodeJsonFromFile(MgoSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
