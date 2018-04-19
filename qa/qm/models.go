package qm

import (
	"errors"
	"sync"

	"github.com/donyori/cqa/common/container"
	"github.com/donyori/cqa/common/container/cmp"
	"github.com/donyori/cqa/data/model"
)

type Candidate struct {
	QuestionId model.Id
	Similarity float32
}

type Request struct {
	Data      *model.Vector32
	TopNumber int
	InC       <-chan *model.QuestionVector
	QuitC     <-chan struct{}
	ExitC     <-chan struct{}
	OutC      chan<- *Candidate
	ErrC      chan<- error
	Wg        *sync.WaitGroup
}

type Response struct {
	Candidates []*Candidate
	Errors     []error
	IsTimeout  bool
}

var (
	ErrReqNilData error = errors.New("data in QM request is nil")
	ErrReqNilInC  error = errors.New("input channel in QM request is nil")
	ErrReqNilOutC error = errors.New("output channel in QM request is nil")
)

func (c *Candidate) Less(another cmp.Comparable) (res bool, err error) {
	a, ok := another.(*Candidate)
	if !ok {
		return false, container.ErrWrongType
	}
	res = false
	if c != nil {
		if a != nil {
			res = c.Similarity < a.Similarity
		}
	} else {
		res = a != nil
	}
	return res, nil
}
