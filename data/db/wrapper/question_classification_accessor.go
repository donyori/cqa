package wrapper

import (
	"errors"

	"github.com/donyori/cqa/data/db/generic"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model"
)

type QuestionClassificationAccessor struct {
	wrappedAccessor
}

var ErrNilQuestionClassificationAccessor error = errors.New(
	"QuestionClassificationAccessor is nil")

func NewQuestionClassificationAccessor(accessor generic.Accessor) (
	qca *QuestionClassificationAccessor, err error) {
	if accessor == nil {
		return nil, generic.ErrNilAccessor
	}
	return &QuestionClassificationAccessor{
		wrappedAccessor: wrappedAccessor{accessor: accessor},
	}, nil
}

func (qca *QuestionClassificationAccessor) IsExisted(params interface{}) (
	res bool, err error) {
	if qca == nil {
		return false, ErrNilQuestionClassificationAccessor
	}
	return qca.accessor.IsExisted(dbid.QuestionClassificationCollection, params)
}

func (qca *QuestionClassificationAccessor) IsExistedById(id model.Id) (
	res bool, err error) {
	if qca == nil {
		return false, ErrNilQuestionClassificationAccessor
	}
	return qca.accessor.IsExistedById(dbid.QuestionClassificationCollection, id)
}

func (qca *QuestionClassificationAccessor) FetchOne(params interface{}) (
	qc *model.QuestionClassification, err error) {
	if qca == nil {
		return nil, ErrNilQuestionClassificationAccessor
	}
	res, err := qca.accessor.FetchOne(
		dbid.QuestionClassificationCollection, params, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	qc, ok := res.(*model.QuestionClassification)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return qc, nil
}

func (qca *QuestionClassificationAccessor) FetchOneById(id model.Id) (
	qc *model.QuestionClassification, err error) {
	if qca == nil {
		return nil, ErrNilQuestionClassificationAccessor
	}
	res, err := qca.accessor.FetchOneById(
		dbid.QuestionClassificationCollection, id, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	qc, ok := res.(*model.QuestionClassification)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return qc, nil
}

func (qca *QuestionClassificationAccessor) FetchAll(params interface{}) (
	qcs []*model.QuestionClassification, err error) {
	if qca == nil {
		return nil, ErrNilQuestionClassificationAccessor
	}
	res, err := qca.accessor.FetchAll(
		dbid.QuestionClassificationCollection, params, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	qcs, ok := res.([]*model.QuestionClassification)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return qcs, nil
}

func (qca *QuestionClassificationAccessor) FetchAllByIds(ids []model.Id) (
	qcs []*model.QuestionClassification, err error) {
	if qca == nil {
		return nil, ErrNilQuestionClassificationAccessor
	}
	res, err := qca.accessor.FetchAllByIds(
		dbid.QuestionClassificationCollection, ids, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	qcs, ok := res.([]*model.QuestionClassification)
	if !ok {
		return nil, ErrResultTypeWrong
	}
	return qcs, nil
}

func (qca *QuestionClassificationAccessor) Scan(
	params interface{}, bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.QuestionClassification,
	resC <-chan error, err error) {
	if qca == nil {
		return nil, nil, ErrNilQuestionClassificationAccessor
	}
	return qca.scan(params, bufferSize, quitC, false)
}

func (qca *QuestionClassificationAccessor) ScanByIds(
	ids []model.Id, bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.QuestionClassification,
	resC <-chan error, err error) {
	if qca == nil {
		return nil, nil, ErrNilQuestionClassificationAccessor
	}
	return qca.scan(ids, bufferSize, quitC, true)
}

func (qca *QuestionClassificationAccessor) Count(params interface{}) (
	res int64, err error) {
	if qca == nil {
		return 0, ErrNilQuestionClassificationAccessor
	}
	return qca.accessor.Count(dbid.QuestionClassificationCollection, params)
}

func (qca *QuestionClassificationAccessor) SaveOne(selector interface{},
	qc *model.QuestionClassification) (isNew bool, err error) {
	if qca == nil {
		return false, ErrNilQuestionClassificationAccessor
	}
	return qca.accessor.SaveOne(
		dbid.QuestionClassificationCollection, selector, qc)
}

func (qca *QuestionClassificationAccessor) SaveOneById(id model.Id,
	qc *model.QuestionClassification) (isNew bool, err error) {
	if qca == nil {
		return false, ErrNilQuestionClassificationAccessor
	}
	return qca.accessor.SaveOneById(
		dbid.QuestionClassificationCollection, id, qc)
}

func (qca *QuestionClassificationAccessor) scan(
	params interface{}, bufferSize uint32,
	quitC <-chan struct{}, paramsAreIds bool) (
	outC <-chan *model.QuestionClassification, resC <-chan error, err error) {
	quitChannel := make(chan struct{}, 1)
	var out <-chan interface{}
	var res <-chan error
	if paramsAreIds {
		out, res, err = qca.accessor.ScanByIds(
			dbid.QuestionClassificationCollection,
			params, bufferSize, quitChannel, nil)
	} else {
		out, res, err = qca.accessor.Scan(dbid.QuestionClassificationCollection,
			params, bufferSize, quitChannel, nil)
	}
	if err != nil {
		return nil, nil, err
	}
	outChannel := make(chan *model.QuestionClassification, cap(out))
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
				qc, ok := q.(*model.QuestionClassification)
				if ok {
					outChannel <- qc
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
					qc, ok := q.(*model.QuestionClassification)
					if ok {
						outChannel <- qc
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
