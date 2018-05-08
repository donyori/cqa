package generic

import (
	"errors"
	"reflect"

	dbid "github.com/donyori/cqa/data/db/id"
)

type Accessor interface {
	WithSession

	FetchOne(cid dbid.CollectionId, params interface{},
		modelType reflect.Type) (res interface{}, err error)
	FetchOneById(cid dbid.CollectionId, id interface{},
		modelType reflect.Type) (res interface{}, err error)
	FetchAll(cid dbid.CollectionId, params interface{},
		modelType reflect.Type) (res interface{}, err error)
	FetchAllByIds(cid dbid.CollectionId, ids interface{},
		modelType reflect.Type) (res interface{}, err error)
	Scan(cid dbid.CollectionId, params interface{}, bufferSize uint32,
		quitC <-chan struct{}, modelType reflect.Type) (
		outC <-chan interface{}, resC <-chan error, err error)
	ScanByIds(cid dbid.CollectionId, ids interface{}, bufferSize uint32,
		quitC <-chan struct{}, modelType reflect.Type) (
		outC <-chan interface{}, resC <-chan error, err error)
	SaveOne(cid dbid.CollectionId, selector interface{}, model interface{}) (
		isNew bool, err error)
	SaveOneById(cid dbid.CollectionId, id interface{}, model interface{}) (
		isNew bool, err error)
}

var (
	ErrNilAccessor error = errors.New("Accessor is nil")
	ErrNilId       error = errors.New("ID is nil")
	ErrEmptyIds    error = errors.New("ID slice is empty")
	ErrIdsNotSlice error = errors.New("IDs must be a slice")
)
