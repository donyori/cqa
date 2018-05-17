package wrapper

import (
	"errors"

	"github.com/donyori/cqa/data/db/generic"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model"
)

type QuestionAccessor struct {
	wrappedAccessor
}

var ErrNilQuestionAccessor error = errors.New("QuestionAccessor is nil")

func NewQuestionAccessor(accessor generic.Accessor) (
	questionAccessor *QuestionAccessor, err error) {
	if accessor == nil {
		return nil, generic.ErrNilAccessor
	}
	return &QuestionAccessor{
		wrappedAccessor: wrappedAccessor{accessor: accessor},
	}, nil
}

func (qa *QuestionAccessor) IsExisted(params interface{}) (
	res bool, err error) {
	if qa == nil {
		return false, ErrNilQuestionAccessor
	}
	return qa.accessor.IsExisted(dbid.QuestionCollection, params)
}

func (qa *QuestionAccessor) IsExistedById(id model.Id) (res bool, err error) {
	if qa == nil {
		return false, ErrNilQuestionAccessor
	}
	return qa.accessor.IsExistedById(dbid.QuestionCollection, id)
}

func (qa *QuestionAccessor) FetchOne(params interface{}) (
	question *model.Question, err error) {
	if qa == nil {
		return nil, ErrNilQuestionAccessor
	}
	res, err := qa.accessor.FetchOne(dbid.QuestionCollection, params, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	question, ok := res.(*model.Question)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return question, nil
}

func (qa *QuestionAccessor) FetchOneById(id model.Id) (
	question *model.Question, err error) {
	if qa == nil {
		return nil, ErrNilQuestionAccessor
	}
	res, err := qa.accessor.FetchOneById(dbid.QuestionCollection, id, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	question, ok := res.(*model.Question)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return question, nil
}

func (qa *QuestionAccessor) FetchAll(params interface{}) (
	questions []*model.Question, err error) {
	if qa == nil {
		return nil, ErrNilQuestionAccessor
	}
	res, err := qa.accessor.FetchAll(dbid.QuestionCollection, params, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	questions, ok := res.([]*model.Question)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return questions, nil
}

func (qa *QuestionAccessor) FetchAllByIds(ids []model.Id) (
	questions []*model.Question, err error) {
	if qa == nil {
		return nil, ErrNilQuestionAccessor
	}
	res, err := qa.accessor.FetchAllByIds(dbid.QuestionCollection, ids, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	questions, ok := res.([]*model.Question)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return questions, nil
}

func (qa *QuestionAccessor) Scan(params interface{}, bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.Question,
	resC <-chan error, err error) {
	if qa == nil {
		return nil, nil, ErrNilQuestionAccessor
	}
	return qa.scan(params, bufferSize, quitC, false)
}

func (qa *QuestionAccessor) ScanByIds(ids []model.Id, bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.Question,
	resC <-chan error, err error) {
	if qa == nil {
		return nil, nil, ErrNilQuestionAccessor
	}
	return qa.scan(ids, bufferSize, quitC, true)
}

func (qa *QuestionAccessor) Count(params interface{}) (res int64, err error) {
	if qa == nil {
		return 0, ErrNilQuestionAccessor
	}
	return qa.accessor.Count(dbid.QuestionCollection, params)
}

func (qa *QuestionAccessor) SaveOne(selector interface{},
	question *model.Question) (isNew bool, err error) {
	if qa == nil {
		return false, ErrNilQuestionAccessor
	}
	return qa.accessor.SaveOne(dbid.QuestionCollection, selector, question)
}

func (qa *QuestionAccessor) SaveOneById(id model.Id,
	question *model.Question) (isNew bool, err error) {
	if qa == nil {
		return false, ErrNilQuestionAccessor
	}
	return qa.accessor.SaveOneById(dbid.QuestionCollection, id, question)
}

func (qa *QuestionAccessor) scan(params interface{}, bufferSize uint32,
	quitC <-chan struct{}, paramsAreIds bool) (outC <-chan *model.Question,
	resC <-chan error, err error) {
	quitChannel := make(chan struct{}, 1)
	var out <-chan interface{}
	var res <-chan error
	if paramsAreIds {
		out, res, err = qa.accessor.ScanByIds(dbid.QuestionCollection,
			params, bufferSize, quitChannel, nil)
	} else {
		out, res, err = qa.accessor.Scan(dbid.QuestionCollection,
			params, bufferSize, quitChannel, nil)
	}
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
