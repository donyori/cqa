package mongodb

import (
	"gopkg.in/mgo.v2"

	"github.com/donyori/cqa/data/dtype"
)

type MgoQuestionAccessor struct {
	MgoConnector
}

func NewMgoQuestionAccessor(settings *MongoDbSettings) *MgoQuestionAccessor {
	return &MgoQuestionAccessor{MgoConnector: *NewMgoConnector(settings)}
}

func (mqa *MgoQuestionAccessor) Get(params interface{}) (
	question *dtype.Question, err error) {
	if params == nil {
		params = NewQueryParams()
	}
	qp, ok := params.(*QueryParams)
	if !ok {
		return nil, ErrNotQueryParams
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

func (mqa *MgoQuestionAccessor) Scan(bufferSize int, params interface{}) (
	out <-chan *dtype.Question, res <-chan error, quit chan<- struct{}, err error) {
	var qp *QueryParams
	if params != nil {
		var ok bool
		qp, ok = params.(*QueryParams)
		if !ok {
			return nil, nil, nil, ErrNotQueryParams
		}
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
	info, err := c.UpsertId(question.QuestionID, question)
	if err != nil {
		return false, err
	}
	return info.Updated == 0, nil
}
