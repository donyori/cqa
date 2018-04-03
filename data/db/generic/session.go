package generic

import (
	"errors"
)

type Session interface {
	IsValid() bool
	GetUrl() (url string, err error)
	Acquire() (nativeSession interface{}, settings interface{}, err error)
	Release()
	Close()
}

type WithSession interface {
	GetSession() Session
	SetSession(session Session) error
}

var ErrNilSession error = errors.New("Session is nil")
