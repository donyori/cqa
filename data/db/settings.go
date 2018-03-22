package db

import (
	"errors"

	"github.com/donyori/cqa/common/json"
)

type DbSettings struct {
	DbType string `json:"db_type"`
}

const (
	DbTypeMongoDB string = "MongoDB"

	DbSettingsFilename string = "settings/db.json"
)

var (
	GlobalSettings DbSettings

	ErrUnknownDbType error = errors.New("db_type is unknown")
)

func init() {
	// Default values:
	GlobalSettings.DbType = DbTypeMongoDB

	err := json.DecodeJsonFromFileIfExist(DbSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}
