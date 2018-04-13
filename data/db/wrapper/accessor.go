package wrapper

import (
	"errors"

	"github.com/donyori/cqa/data/db/generic"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model"
)

type wrappedAccessor struct {
	accessor generic.Accessor
}

type QuestionAccessor struct {
	wrappedAccessor
}

type QuestionVectorAccessor struct {
	wrappedAccessor
}

var (
	ErrNilQuestionAccessor       error = errors.New("QuestionAccessor is nil")
	ErrNilQuestionVectorAccessor error = errors.New(
		"QuestionVectorAccessor is nil")
	ErrResultTypeWrong error = errors.New("result type is wrong")
)

func (wa *wrappedAccessor) GetAccessor() generic.Accessor {
	if wa == nil {
		return nil
	}
	return wa.accessor
}

func NewQuestionAccessor(accessor generic.Accessor) (
	questionAccessor *QuestionAccessor, err error) {
	if accessor == nil {
		return nil, generic.ErrNilAccessor
	}
	return &QuestionAccessor{
		wrappedAccessor: wrappedAccessor{accessor: accessor},
	}, nil
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

func NewQuestionVectorAccessor(accessor generic.Accessor) (
	qva *QuestionVectorAccessor, err error) {
	if accessor == nil {
		return nil, generic.ErrNilAccessor
	}
	return &QuestionVectorAccessor{
		wrappedAccessor: wrappedAccessor{accessor: accessor},
	}, nil
}

func (qva *QuestionVectorAccessor) FetchOne(params interface{}) (
	qv *model.QuestionVector, err error) {
	if qva == nil {
		return nil, ErrNilQuestionVectorAccessor
	}
	res, err := qva.accessor.FetchOne(
		dbid.QuestionVectorCollection, params, nil)
	if err != nil {
		return nil, err
	}
	qv, ok := res.(*model.QuestionVector)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return qv, nil
}

func (qva *QuestionVectorAccessor) FetchOneById(id model.Id) (
	qv *model.QuestionVector, err error) {
	if qva == nil {
		return nil, ErrNilQuestionVectorAccessor
	}
	res, err := qva.accessor.FetchOneById(
		dbid.QuestionVectorCollection, id, nil)
	if err != nil {
		return nil, err
	}
	qv, ok := res.(*model.QuestionVector)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return qv, nil
}

func (qva *QuestionVectorAccessor) FetchAll(params interface{}) (
	qvs []*model.QuestionVector, err error) {
	if qva == nil {
		return nil, ErrNilQuestionVectorAccessor
	}
	res, err := qva.accessor.FetchAll(
		dbid.QuestionVectorCollection, params, nil)
	if err != nil {
		return nil, err
	}
	qvs, ok := res.([]*model.QuestionVector)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return qvs, nil
}

func (qva *QuestionVectorAccessor) FetchAllByIds(ids []model.Id) (
	qvs []*model.QuestionVector, err error) {
	if qva == nil {
		return nil, ErrNilQuestionVectorAccessor
	}
	res, err := qva.accessor.FetchAllByIds(
		dbid.QuestionVectorCollection, ids, nil)
	if err != nil {
		return nil, err
	}
	qvs, ok := res.([]*model.QuestionVector)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return qvs, nil
}

func (qva *QuestionVectorAccessor) Scan(params interface{}, bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.QuestionVector,
	resC <-chan error, err error) {
	if qva == nil {
		return nil, nil, ErrNilQuestionVectorAccessor
	}
	return qva.scan(params, bufferSize, quitC, false)
}

func (qva *QuestionVectorAccessor) ScanByIds(ids []model.Id, bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.QuestionVector,
	resC <-chan error, err error) {
	if qva == nil {
		return nil, nil, ErrNilQuestionVectorAccessor
	}
	return qva.scan(ids, bufferSize, quitC, true)
}

func (qva *QuestionVectorAccessor) SaveOne(selector interface{},
	qv *model.QuestionVector) (isNew bool, err error) {
	if qva == nil {
		return false, ErrNilQuestionVectorAccessor
	}
	return qva.accessor.SaveOne(
		dbid.QuestionVectorCollection, selector, qv)
}

func (qva *QuestionVectorAccessor) SaveOneById(id model.Id,
	qv *model.QuestionVector) (isNew bool, err error) {
	if qva == nil {
		return false, ErrNilQuestionVectorAccessor
	}
	return qva.accessor.SaveOneById(dbid.QuestionVectorCollection, id, qv)
}

func (qva *QuestionVectorAccessor) scan(params interface{}, bufferSize uint32,
	quitC <-chan struct{}, paramsAreIds bool) (
	outC <-chan *model.QuestionVector, resC <-chan error, err error) {
	quitChannel := make(chan struct{}, 1)
	var out <-chan interface{}
	var res <-chan error
	if paramsAreIds {
		out, res, err = qva.accessor.ScanByIds(dbid.QuestionVectorCollection,
			params, bufferSize, quitChannel, nil)
	} else {
		out, res, err = qva.accessor.Scan(dbid.QuestionVectorCollection,
			params, bufferSize, quitChannel, nil)
	}
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
				qv, ok := q.(*model.QuestionVector)
				if ok {
					outChannel <- qv
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
					qv, ok := q.(*model.QuestionVector)
					if ok {
						outChannel <- qv
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
