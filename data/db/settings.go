package db

import (
	"errors"

	"github.com/donyori/cqa/common/json"
)

type DbType int8

type DbSettings struct {
	Type DbType `json:"type"`
}

const (
	DbTypeMongoDB DbType = iota
)

const DbSettingsFilename string = "settings/db.json"

var (
	GlobalSettings DbSettings

	ErrInvalidDbType error = errors.New("DB type is invalid")

	dbTypeStrings = [...]string{
		"MongoDB",
	}
)

func init() {
	// Default values:
	GlobalSettings.Type = DbTypeMongoDB

	_, err := json.DecodeJsonFromFileIfExist(
		DbSettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}

func (dt DbType) String() string {
	if !dt.IsValid() {
		return "Unknown"
	}
	return dbTypeStrings[dt]
}

func (dt DbType) IsValid() bool {
	return dt == DbTypeMongoDB
}
