package qm

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/donyori/cqa/common/container"
	"github.com/donyori/cqa/common/container/topkbuf"
	"github.com/donyori/cqa/nlp"
)

type Requester struct {
	provider *Provider
	matcher  *Matcher

	lock     sync.RWMutex
	initOnce sync.Once
}

var (
	ErrNilRequester         error = errors.New("requester is nil")
	ErrRequesterAlreadyInit error = errors.New(
		"requester has already initialized")
	ErrRequesterNotInit     error = errors.New("requester is NOT initialized")
	ErrNonPositiveTopNumber error = errors.New(
		"top number is non-positive and no default top number can use")
)

func NewRequester(provider *Provider, matcher *Matcher) (
	requester *Requester, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			requester = nil
			e, ok := panicErr.(error)
			if ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicErr)
			}
		}
	}()
	r := new(Requester)
	err = r.Init(provider, matcher)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Requester) Init(provider *Provider, matcher *Matcher) error {
	if r == nil {
		return ErrNilRequester
	}
	if provider == nil {
		return ErrNilProvider
	}
	if matcher == nil {
		return ErrNilMatcher
	}
	isInitialized := true
	r.initOnce.Do(func() {
		isInitialized = false
		r.lock.Lock()
		defer r.lock.Unlock()
		r.provider = provider
		r.matcher = matcher
	})
	if isInitialized {
		r.lock.RLock()
		defer r.lock.RUnlock()
		if provider != r.provider || matcher != r.matcher {
			return ErrRequesterAlreadyInit
		}
	}
	return nil
}

func (r *Requester) Match(question string, topNumber int,
	timeLimit time.Duration) (respC <-chan *Response, err error) {
	if r == nil {
		return nil, ErrNilRequester
	}
	r.lock.RLock()
	defer r.lock.RUnlock()
	if r.provider == nil || r.matcher == nil {
		return nil, ErrRequesterNotInit
	}
	if topNumber <= 0 {
		topNumber = GlobalSettings.DefaultTopNumber
		if topNumber <= 0 {
			return nil, ErrNonPositiveTopNumber
		}
	}
	if timeLimit < 0 {
		timeLimit = GlobalSettings.GetDefaultTimeLimit()
	}
	var timer *time.Timer
	if timeLimit > 0 {
		timer = time.NewTimer(timeLimit)
	}
	candidateBuffer, err := topkbuf.NewTopKBuffer(topNumber)
	if err != nil {
		return nil, err
	}
	vector, err := nlp.Embedding(question)
	if err != nil {
		return nil, err
	}
	inQuitC := make(chan struct{})
	inC, err := r.provider.ProvideQuestionVector(
		GlobalSettings.InputBufferSize, inQuitC)
	if err != nil {
		return nil, err
	}
	quitC := make(chan struct{})
	exitC := make(chan struct{})
	outC := make(chan *Candidate, GlobalSettings.OutputBufferSize)
	errC := make(chan error, GlobalSettings.ErrorBufferSize)
	wg := new(sync.WaitGroup)
	req := &Request{
		Data:      vector,
		TopNumber: topNumber,
		InC:       inC,
		QuitC:     quitC,
		ExitC:     exitC,
		OutC:      outC,
		ErrC:      errC,
		Wg:        wg,
	}
	doneC := make(chan struct{})
	var timeOutC chan bool
	wg.Add(1) // For the broadcast goroutine.
	go r.requestChannelsPostProcessing(req, doneC)
	go r.sendQuitSignal(wg, inQuitC)
	if timer != nil {
		timeOutC = make(chan bool, 1)
		go r.timerProcess(question, timer, timeOutC, inQuitC, quitC)
	}
	err = r.matcher.SendRequest(req)
	if err != nil {
		// Matcher didn't dispatch request.
		// Call wg.Done() to stop other goroutines.
		wg.Done()
		<-doneC
		if timeOutC != nil {
			<-timeOutC
		}
		return nil, err
	}
	respChannel := make(chan *Response, 1)
	go r.response(outC, errC, timeOutC, candidateBuffer, respChannel)
	return respChannel, nil
}

func (r *Requester) requestChannelsPostProcessing(req *Request,
	doneC chan<- struct{}) {
	if req.Wg == nil {
		panic(errors.New("wait group is nil"))
	}
	if doneC != nil {
		defer close(doneC)
	}
	defer func() {
		if req.InC != nil {
			// Drain req.InC.
			for _ = range req.InC {
			}
		}
	}()
	defer func() {
		if req.ErrC != nil {
			close(req.ErrC)
		}
	}()
	defer func() {
		if req.OutC != nil {
			close(req.OutC)
		}
	}()
	req.Wg.Wait()
	if req.InC != nil && req.ErrC != nil {
		rin := len(req.InC)
		if rin > 0 {
			req.ErrC <- fmt.Errorf("remaining inputs number: %d", rin)
		}
	}
}

func (r *Requester) sendQuitSignal(wg *sync.WaitGroup,
	quitC chan<- struct{}) {
	if wg == nil {
		panic(errors.New("wait group is nil"))
	}
	if quitC == nil {
		panic(errors.New("quit channel is nil"))
	}
	defer close(quitC)
	wg.Wait()
}

func (r *Requester) timerProcess(q string, timer *time.Timer,
	outC chan<- bool, quitC <-chan struct{}, reqQuitC chan<- struct{}) {
	isTimeOut := false
	if outC != nil {
		defer func() {
			outC <- isTimeOut
			close(outC)
		}()
	}
	if reqQuitC == nil {
		panic(errors.New("request quit channel is nil"))
	}
	if timer == nil {
		return
	}
	select {
	case <-quitC:
		if !timer.Stop() {
			// Drain time.C
			<-timer.C
		}
	case <-timer.C:
		defer log.Printf("Question matching time out: %q\n", q)
		isTimeOut = true
		close(reqQuitC)
	}
}

func (r *Requester) response(candidateC <-chan *Candidate, errC <-chan error,
	isTimeOutC <-chan bool, candidateBuffer *topkbuf.TopKBuffer,
	respC chan<- *Response) {
	if respC == nil {
		panic(errors.New("response channel is nil"))
	}
	defer close(respC)
	var candidates []*Candidate
	var errs []error
	defer func() {
		isTimeOut := false
		if isTimeOutC != nil {
			isTimeOut = <-isTimeOutC
		}
		resp := &Response{
			Candidates: candidates,
			Errors:     errs,
			IsTimeOut:  isTimeOut,
		}
		respC <- resp
	}()
	if candidateBuffer == nil {
		err := errors.New("candidate buffer is nil")
		errs = append(errs, err)
		panic(err)
	}
	for candidateC != nil || errC != nil {
		select {
		case candidate, ok := <-candidateC:
			if ok {
				if candidate != nil {
					candidateBuffer.Add(candidate)
				}
			} else {
				candidateC = nil
			}
		case err, ok := <-errC:
			if ok {
				if err != nil {
					errs = append(errs, err)
				}
			} else {
				errC = nil
			}
		}
	}
	outs, err := candidateBuffer.Flush()
	if err != nil {
		errs = append(errs, err)
		panic(err)
	}
	for _, out := range outs {
		candidate, ok := out.(*Candidate)
		if !ok {
			errs = append(errs, container.ErrWrongType)
		}
		candidates = append(candidates, candidate)
	}
}
