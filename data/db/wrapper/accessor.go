package wrapper

import (
	"errors"

	"github.com/donyori/cqa/data/db/generic"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model"
	"github.com/donyori/cqa/data/model/helper"
)

type QuestionAccessor struct {
	accessor generic.Accessor
}

type QuestionVectorAccessor struct {
	accessor generic.Accessor
}

var (
	ErrNilQuestionAccessor       error = errors.New("QuestionAccessor is nil")
	ErrNilQuestionVectorAccessor error = errors.New(
		"QuestionVectorAccessor is nil")
	ErrResultTypeWrong error = errors.New("result type is wrong")
)

func NewQuestionAccessor(accessor generic.Accessor) (
	questionAccessor *QuestionAccessor, err error) {
	if accessor == nil {
		return nil, generic.ErrNilAccessor
	}
	return &QuestionAccessor{accessor: accessor}, nil
}

func (qa *QuestionAccessor) Get(params interface{}) (
	question *model.Question, err error) {
	if qa == nil {
		return nil, ErrNilQuestionAccessor
	}
	res, err := qa.accessor.Get(dbid.QuestionCollection,
		params, helper.QuestionMaker)
	if err != nil {
		return nil, err
	}
	question, ok := res.(*model.Question)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return question, nil
}

func (qa *QuestionAccessor) GetById(id interface{}) (
	question *model.Question, err error) {
	if qa == nil {
		return nil, ErrNilQuestionAccessor
	}
	res, err := qa.accessor.GetById(dbid.QuestionCollection,
		id, helper.QuestionMaker)
	if err != nil {
		return nil, err
	}
	question, ok := res.(*model.Question)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return question, nil
}

func (qa *QuestionAccessor) Scan(params interface{}, bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.Question,
	resC <-chan error, err error) {
	if qa == nil {
		return nil, nil, ErrNilQuestionAccessor
	}
	quitChannel := make(chan struct{}, 1)
	out, res, err := qa.accessor.Scan(dbid.QuestionCollection,
		params, bufferSize, quitChannel, helper.QuestionMaker)
	if err != nil {
		return nil, nil, err
	}
	outChannel := make(chan *model.Question, cap(out))
	resultChannel := make(chan error, 1)
	go func() {
		defer close(resultChannel)
		defer close(outChannel)
		defer func() {
			if quitChannel != nil {
				close(quitChannel)
			}
		}()
		var e error
		isQuit := false
		for !isQuit {
			select {
			case msg, ok := <-quitC:
				if quitChannel == nil {
					break
				}
				if ok {
					quitChannel <- msg
				} else {
					close(quitChannel)
					quitChannel = nil
				}
			case q, ok := <-out:
				if !ok {
					e = <-res
					isQuit = true
					break
				}
				question, ok := q.(*model.Question)
				if ok {
					outChannel <- question
				} else {
					if quitChannel != nil {
						quitChannel <- struct{}{}
					}
					e = ErrResultTypeWrong
					isQuit = true
				}
			case e = <-res:
				// Drain channel
				for q := range out {
					question, ok := q.(*model.Question)
					if ok {
						outChannel <- question
					} else {
						e = ErrResultTypeWrong
						break
					}
				}
				isQuit = true
			}
		}
		resultChannel <- e
	}()
	return outChannel, resultChannel, nil
}

func (qa *QuestionAccessor) Save(selector interface{},
	question *model.Question) (isNew bool, err error) {
	if qa == nil {
		return false, ErrNilQuestionAccessor
	}
	return qa.accessor.Save(dbid.QuestionCollection, selector, question)
}

func (qa *QuestionAccessor) SaveById(id interface{},
	question *model.Question) (isNew bool, err error) {
	if qa == nil {
		return false, ErrNilQuestionAccessor
	}
	return qa.accessor.SaveById(dbid.QuestionCollection, id, question)
}

func NewQuestionVectorAccessor(accessor generic.Accessor) (
	questionAccessor *QuestionVectorAccessor, err error) {
	if accessor == nil {
		return nil, generic.ErrNilAccessor
	}
	return &QuestionVectorAccessor{accessor: accessor}, nil
}

func (qva *QuestionVectorAccessor) Get(params interface{}) (
	question *model.QuestionVector, err error) {
	if qva == nil {
		return nil, ErrNilQuestionVectorAccessor
	}
	res, err := qva.accessor.Get(dbid.QuestionVectorCollection,
		params, helper.QuestionVectorMaker)
	if err != nil {
		return nil, err
	}
	question, ok := res.(*model.QuestionVector)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return question, nil
}

func (qva *QuestionVectorAccessor) GetById(id interface{}) (
	question *model.QuestionVector, err error) {
	if qva == nil {
		return nil, ErrNilQuestionVectorAccessor
	}
	res, err := qva.accessor.GetById(dbid.QuestionVectorCollection,
		id, helper.QuestionVectorMaker)
	if err != nil {
		return nil, err
	}
	question, ok := res.(*model.QuestionVector)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return question, nil
}

func (qva *QuestionVectorAccessor) Scan(params interface{}, bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.QuestionVector,
	resC <-chan error, err error) {
	if qva == nil {
		return nil, nil, ErrNilQuestionVectorAccessor
	}
	quitChannel := make(chan struct{}, 1)
	out, res, err := qva.accessor.Scan(dbid.QuestionVectorCollection,
		params, bufferSize, quitChannel, helper.QuestionVectorMaker)
	if err != nil {
		return nil, nil, err
	}
	outChannel := make(chan *model.QuestionVector, cap(out))
	resultChannel := make(chan error, 1)
	go func() {
		defer close(resultChannel)
		defer close(outChannel)
		defer func() {
			if quitChannel != nil {
				close(quitChannel)
			}
		}()
		var e error
		isQuit := false
		for !isQuit {
			select {
			case msg, ok := <-quitC:
				if quitChannel == nil {
					break
				}
				if ok {
					quitChannel <- msg
				} else {
					close(quitChannel)
					quitChannel = nil
				}
			case q, ok := <-out:
				if !ok {
					e = <-res
					isQuit = true
					break
				}
				question, ok := q.(*model.QuestionVector)
				if ok {
					outChannel <- question
				} else {
					if quitChannel != nil {
						quitChannel <- struct{}{}
					}
					e = ErrResultTypeWrong
					isQuit = true
				}
			case e = <-res:
				// Drain channel
				for q := range out {
					question, ok := q.(*model.QuestionVector)
					if ok {
						outChannel <- question
					} else {
						e = ErrResultTypeWrong
						break
					}
				}
				isQuit = true
			}
		}
		resultChannel <- e
	}()
	return outChannel, resultChannel, nil
}

func (qva *QuestionVectorAccessor) Save(selector interface{},
	question *model.QuestionVector) (isNew bool, err error) {
	if qva == nil {
		return false, ErrNilQuestionVectorAccessor
	}
	return qva.accessor.Save(dbid.QuestionVectorCollection, selector, question)
}

func (qva *QuestionVectorAccessor) SaveById(id interface{},
	question *model.QuestionVector) (isNew bool, err error) {
	if qva == nil {
		return false, ErrNilQuestionVectorAccessor
	}
	return qva.accessor.SaveById(dbid.QuestionVectorCollection, id, question)
}
