package generic

import (
	"errors"

	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model/helper"
)

type Accessor interface {
	WithSession

	Get(cid dbid.CollectionId, params interface{}, maker helper.Maker) (
		res interface{}, err error)
	GetById(cid dbid.CollectionId, id interface{}, maker helper.Maker) (
		res interface{}, err error)
	Scan(cid dbid.CollectionId, params interface{}, bufferSize int,
		quitC <-chan struct{}, maker helper.Maker) (
		outC <-chan interface{}, resC <-chan error, err error)
	Save(cid dbid.CollectionId, selector interface{}, model interface{}) (
		isNew bool, err error)
	SaveById(cid dbid.CollectionId, id interface{}, model interface{}) (
		isNew bool, err error)
}

var (
	ErrNilAccessor error = errors.New("Accessor is nil")
	ErrNilId       error = errors.New("ID is nil")
)
