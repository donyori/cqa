package db

import (
	"errors"

	"github.com/donyori/cqa/common/json"
)

type DbType int8

type Settings struct {
	Type DbType `json:"type"`
}

const (
	DbTypeMongoDB DbType = iota
)

const SettingsFilename string = "settings/db.json"

var (
	GlobalSettings Settings

	ErrInvalidDbType error = errors.New("DB type is invalid")

	dbTypeStrings = [...]string{
		"MongoDB",
	}
)

func init() {
	// Default values:
	GlobalSettings.Type = DbTypeMongoDB

	_, err := json.DecodeJsonFromFileIfExist(
		SettingsFilename, &GlobalSettings)
	if err != nil {
		panic(err)
	}
}

func (dt DbType) IsValid() bool {
	return dt == DbTypeMongoDB
}

func (dt DbType) String() string {
	if !dt.IsValid() {
		return "Unknown"
	}
	return dbTypeStrings[dt]
}
