package db

import (
	"github.com/donyori/cqa/data/db/generic"
	"github.com/donyori/cqa/data/db/mongodb"
)

func NewSession() (session generic.Session, err error) {
	dbType := GlobalSettings.Type
	switch dbType {
	case DbTypeMongoDB:
		return mongodb.NewMgoSession(nil)
	default:
		return nil, ErrInvalidDbType
	}
}

func CleanUpSessionPool() {
	dbType := GlobalSettings.Type
	switch dbType {
	case DbTypeMongoDB:
		mongodb.CleanUpSessionPool()
	}
}
