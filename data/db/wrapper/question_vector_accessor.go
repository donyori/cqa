package wrapper

import (
	"errors"

	"github.com/donyori/cqa/data/db/generic"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model"
)

type QuestionVectorAccessor struct {
	wrappedAccessor
}

var ErrNilQuestionVectorAccessor error = errors.New(
	"QuestionVectorAccessor is nil")

func NewQuestionVectorAccessor(accessor generic.Accessor) (
	qva *QuestionVectorAccessor, err error) {
	if accessor == nil {
		return nil, generic.ErrNilAccessor
	}
	return &QuestionVectorAccessor{
		wrappedAccessor: wrappedAccessor{accessor: accessor},
	}, nil
}

func (qva *QuestionVectorAccessor) IsExisted(params interface{}) (
	res bool, err error) {
	if qva == nil {
		return false, ErrNilQuestionVectorAccessor
	}
	return qva.accessor.IsExisted(dbid.QuestionVectorCollection, params)
}

func (qva *QuestionVectorAccessor) IsExistedById(id model.Id) (
	res bool, err error) {
	if qva == nil {
		return false, ErrNilQuestionVectorAccessor
	}
	return qva.accessor.IsExistedById(dbid.QuestionVectorCollection, id)
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

func (qva *QuestionVectorAccessor) Count(params interface{}) (
	res int64, err error) {
	if qva == nil {
		return 0, ErrNilQuestionVectorAccessor
	}
	return qva.accessor.Count(dbid.QuestionVectorCollection, params)
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
