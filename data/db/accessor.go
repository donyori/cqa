package db

import (
	"github.com/donyori/cqa/data/db/generic"
	"github.com/donyori/cqa/data/db/mongodb"
)

func NewAccessor(session generic.Session) (
	accessor generic.Accessor, err error) {
	dbType := GlobalSettings.Type
	switch dbType {
	case DbTypeMongoDB:
		return mongodb.NewMgoAccessor(session)
	default:
		return nil, ErrInvalidDbType
	}
}
