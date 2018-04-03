package generic

import (
	"errors"

	"github.com/donyori/cqa/data/db/id"
)

type Maintainer interface {
	WithSession

	EnsureIndexes(cid id.CollectionId, isBackground bool) error
	EnsureDataTypes(cid id.CollectionId) error
}

var ErrNilMaintainer error = errors.New("Maintainer is nil")
