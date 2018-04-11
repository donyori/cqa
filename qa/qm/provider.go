package qm

import (
	"errors"
	"fmt"
	"sync"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/wrapper"
	"github.com/donyori/cqa/data/model"
)

type Provider struct {
	qvBuf []*model.QuestionVector

	initOnce sync.Once
}

var ErrNilProvider error = errors.New("provider is nil")

func NewProvider() (provider *Provider, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			provider = nil
			e, ok := panicErr.(error)
			if ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicErr)
			}
		}
	}()
	p := new(Provider)
	p.Init()
	return p, nil
}

func (p *Provider) Init() {
	if p == nil {
		panic(ErrNilProvider)
	}
	p.initOnce.Do(func() {
		if GlobalSettings.EnableQuestionVectorBuffer {
			err := p.bufferQuestionVectors()
			if err != nil {
				panic(err)
			}
		} else {
			p.qvBuf = nil
		}
	})
}

func (p *Provider) ProvideQuestionVector(bufferSize uint32,
	quitC <-chan struct{}) (outC <-chan *model.QuestionVector, err error) {
	if p == nil {
		return nil, ErrNilProvider
	}
	p.Init()
	outChannel := make(chan *model.QuestionVector, bufferSize)
	if p.qvBuf != nil {
		go p.provideQuestionVectorFromBuffer(outChannel, quitC)
	} else {
		sess, err := db.NewSession()
		if err != nil {
			return nil, err
		}
		// defer sess.Close() in other goroutine.
		accessor, err := db.NewAccessor(sess)
		if err != nil {
			sess.Close()
			return nil, err
		}
		qva, err := wrapper.NewQuestionVectorAccessor(accessor)
		if err != nil {
			sess.Close()
			return nil, err
		}
		go p.provideQuestionVectorFromDb(qva, outChannel, quitC)
	}
	return outChannel, nil
}

func (p *Provider) bufferQuestionVectors() error {
	if p == nil {
		return ErrNilMatcher
	}
	sess, err := db.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	accessor, err := db.NewAccessor(sess)
	if err != nil {
		return err
	}
	qva, err := wrapper.NewQuestionVectorAccessor(accessor)
	if err != nil {
		return err
	}
	outC, resC, err := qva.Scan(nil, 4, nil)
	if err != nil {
		return err
	}
	p.qvBuf = nil
	for qv := range outC {
		p.qvBuf = append(p.qvBuf, qv)
	}
	return <-resC
}

func (p *Provider) provideQuestionVectorFromBuffer(
	outC chan<- *model.QuestionVector, quitC <-chan struct{}) {
	if outC == nil {
		panic(errors.New("out channel is nil"))
	}
	defer close(outC)
	if p == nil {
		panic(ErrNilProvider)
	}
	if p.qvBuf == nil {
		return
	}
	for _, qv := range p.qvBuf {
		select {
		case <-quitC:
			return
		default:
			outC <- qv
		}
	}
}

func (p *Provider) provideQuestionVectorFromDb(
	qva *wrapper.QuestionVectorAccessor,
	outC chan<- *model.QuestionVector, quitC <-chan struct{}) {
	if outC == nil {
		panic(errors.New("out channel is nil"))
	}
	defer close(outC)
	if qva == nil {
		panic(errors.New("question vector accessor is nil"))
	}
	sess := qva.GetAccessor().GetSession()
	defer sess.Close()
	qvqc := make(chan struct{})
	defer close(qvqc)
	qvc, _, err := qva.Scan(nil, uint32(cap(outC)), qvqc)
	if err != nil {
		panic(err)
	}
	isDone := false
	for !isDone {
		select {
		case <-quitC:
			isDone = true
		case qv, ok := <-qvc:
			if ok {
				outC <- qv
			} else {
				isDone = true
			}
		}
	}
}
