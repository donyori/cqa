package mongodb

import (
	"gopkg.in/mgo.v2"

	"github.com/donyori/cqa/data/dtype"
)

type MgoQuestionAccessor struct {
	MgoConnector
}

type MgoQuestionVectorAccessor struct {
	MgoConnector
}

func NewMgoQuestionAccessor(settings *MongoDbSettings) *MgoQuestionAccessor {
	return &MgoQuestionAccessor{MgoConnector: *NewMgoConnector(settings)}
}

func (mqa *MgoQuestionAccessor) Get(params interface{}) (
	question *dtype.Question, err error) {
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return nil, err
	}
	if qp == nil {
		qp = NewQueryParams()
	}
	mqa.RLock()
	defer mqa.RUnlock()
	if !mqa.isConnectedWithoutLock() {
		return nil, ErrNotConnected
	}
	settings := mqa.getSettings()
	c := mqa.Session.DB(settings.DbName).C(settings.CNames[MgoCNameKeyQ])
	qp.Limit = 1
	q := qp.MakeQuery(c)
	res := new(dtype.Question)
	err = q.One(res)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (mqa *MgoQuestionAccessor) GetById(id interface{}) (
	question *dtype.Question, err error) {
	params := NewQueryParams()
	params.Id = id
	return mqa.Get(params)
}

func (mqa *MgoQuestionAccessor) Scan(params interface{}, bufferSize int) (
	out <-chan *dtype.Question, res <-chan error, quit chan<- struct{}, err error) {
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return nil, nil, nil, err
	}
	mqa.RLock()
	defer mqa.RUnlock()
	if !mqa.isConnectedWithoutLock() {
		return nil, nil, nil, ErrNotConnected
	}
	var outChan chan *dtype.Question
	if bufferSize > 0 {
		outChan = make(chan *dtype.Question, bufferSize)
	} else {
		outChan = make(chan *dtype.Question)
	}
	resChan := make(chan error, 1)
	quitChan := make(chan struct{}, 1)
	go func() {
		mqa.RLock()
		defer mqa.RUnlock()
		settings := mqa.getSettings()
		c := mqa.Session.DB(settings.DbName).C(settings.CNames[MgoCNameKeyQ])
		q := qp.MakeQuery(c)
		iter := q.Iter()
		defer iter.Close() // Ignore error.
		defer close(resChan)
		defer close(outChan)
		result := new(dtype.Question)
		isQuit := false
		for !isQuit && iter.Next(result) {
			select {
			case <-quitChan:
				isQuit = true
			default:
				outChan <- result
				result = new(dtype.Question) // Make a new struct each time.
			}
		}
		iterErr := iter.Err()
		resChan <- iterErr
	}()
	return outChan, resChan, quitChan, nil
}

func (mqa *MgoQuestionAccessor) Save(question *dtype.Question) (isNew bool, err error) {
	mqa.RLock()
	defer mqa.RUnlock()
	if !mqa.isConnectedWithoutLock() {
		return false, ErrNotConnected
	}
	settings := mqa.getSettings()
	c := mqa.Session.DB(settings.DbName).C(settings.CNames[MgoCNameKeyQ])
	info, err := c.UpsertId(question.QuestionId, question)
	if err != nil {
		return false, err
	}
	return info.Updated == 0, nil
}

func NewMgoQuestionVectorAccessor(settings *MongoDbSettings) *MgoQuestionVectorAccessor {
	return &MgoQuestionVectorAccessor{MgoConnector: *NewMgoConnector(settings)}
}

func (mqva *MgoQuestionVectorAccessor) Get(params interface{}) (
	questionVector *dtype.QuestionVector, err error) {
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return nil, err
	}
	if qp == nil {
		qp = NewQueryParams()
	}
	mqva.RLock()
	defer mqva.RUnlock()
	if !mqva.isConnectedWithoutLock() {
		return nil, ErrNotConnected
	}
	settings := mqva.getSettings()
	c := mqva.Session.DB(settings.DbName).C(settings.CNames[MgoCNameKeyQv])
	qp.Limit = 1
	q := qp.MakeQuery(c)
	res := new(dtype.QuestionVector)
	err = q.One(res)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (mqva *MgoQuestionVectorAccessor) GetById(id interface{}) (
	questionVector *dtype.QuestionVector, err error) {
	params := NewQueryParams()
	params.Id = id
	return mqva.Get(params)
}

func (mqva *MgoQuestionVectorAccessor) Scan(params interface{}, bufferSize int) (
	out <-chan *dtype.QuestionVector, res <-chan error, quit chan<- struct{}, err error) {
	qp, err := ConvertToQueryParams(params)
	if err != nil {
		return nil, nil, nil, err
	}
	mqva.RLock()
	defer mqva.RUnlock()
	if mqva.isConnectedWithoutLock() {
		return nil, nil, nil, ErrNotConnected
	}
	var outChan chan *dtype.QuestionVector
	if bufferSize > 0 {
		outChan = make(chan *dtype.QuestionVector, bufferSize)
	} else {
		outChan = make(chan *dtype.QuestionVector)
	}
	resChan := make(chan error, 1)
	quitChan := make(chan struct{}, 1)
	go func() {
		mqva.RLock()
		defer mqva.RUnlock()
		settings := mqva.getSettings()
		c := mqva.Session.DB(settings.DbName).C(settings.CNames[MgoCNameKeyQv])
		q := qp.MakeQuery(c)
		iter := q.Iter()
		defer iter.Close() // Ignore error.
		defer close(resChan)
		defer close(outChan)
		result := new(dtype.QuestionVector)
		isQuit := false
		for !isQuit && iter.Next(result) {
			select {
			case <-quitChan:
				isQuit = true
			default:
				outChan <- result
				result = new(dtype.QuestionVector) // Make a new struct each time.
			}
		}
		iterErr := iter.Err()
		resChan <- iterErr
	}()
	return outChan, resChan, quitChan, nil
}

func (mqva *MgoQuestionVectorAccessor) Save(questionVector *dtype.QuestionVector) (
	isNew bool, err error) {
	mqva.RLock()
	defer mqva.RUnlock()
	if !mqva.isConnectedWithoutLock() {
		return false, ErrNotConnected
	}
	settings := mqva.getSettings()
	c := mqva.Session.DB(settings.DbName).C(settings.CNames[MgoCNameKeyQv])
	info, err := c.UpsertId(questionVector.QuestionId, questionVector)
	if err != nil {
		return false, err
	}
	return info.Updated == 0, nil
}
