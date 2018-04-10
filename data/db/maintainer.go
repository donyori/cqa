package db

import (
	"github.com/donyori/cqa/data/db/generic"
	"github.com/donyori/cqa/data/db/mongodb"
)

func NewMaintainer(session generic.Session) (
	mantainer generic.Maintainer, err error) {
	dbType := GlobalSettings.Type
	switch dbType {
	case DbTypeMongoDB:
		return mongodb.NewMaintainer(session)
	default:
		return nil, ErrInvalidDbType
	}
}
