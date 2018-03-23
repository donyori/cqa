package db

import (
	"github.com/donyori/cqa/data/db/mongodb"
)

type Maintainer interface {
	Connector

	EnsureIndexes(isBackground bool) error
}

func NewMaintainer() (mantainer Maintainer, err error) {
	switch GlobalSettings.DbType {
	case DbTypeMongoDB:
		return mongodb.NewMgoMaintainer(nil), nil
	default:
		return nil, ErrUnknownDbType
	}
}
