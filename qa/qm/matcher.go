package qm

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/donyori/cqa/common/container"
	"github.com/donyori/cqa/common/container/topkbuf"
)

type Matcher struct {
	matcherNumber int

	requestC        chan<- *Request
	quitC           chan<- struct{}
	exitC           chan<- struct{}
	wg              sync.WaitGroup
	logQuitErrDoneC <-chan struct{}

	lock     sync.RWMutex
	initOnce sync.Once
	exitOnce sync.Once
}

var (
	ErrNilMatcher               error = errors.New("matcher is nil")
	ErrNonPositiveMatcherNumber error = errors.New(
		"matcher number is non-positive and no default matcher number can use")
	ErrMatcherAlreadyInit error = errors.New("matcher has already initialized")
	ErrMatcherShutdown    error = errors.New("question matcher shutdown")
	ErrNilRequest         error = errors.New("request is nil")
)

func NewMatcher(matcherNumber int) (matcher *Matcher, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			matcher = nil
			e, ok := panicErr.(error)
			if ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicErr)
			}
		}
	}()
	m := new(Matcher)
	err = m.Init(matcherNumber)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Matcher) Init(matcherNumber int) error {
	if m == nil {
		return ErrNilMatcher
	}
	mn := matcherNumber
	if mn <= 0 {
		mn = GlobalSettings.DefaultMatcherNumber
		if mn <= 0 {
			return ErrNonPositiveMatcherNumber
		}
	}
	isInitialized := true
	m.initOnce.Do(func() {
		isInitialized = false
		m.lock.Lock()
		defer m.lock.Unlock()
		rbs := GlobalSettings.RequestBufferSize
		ebs := GlobalSettings.ErrorBufferSize

		m.matcherNumber = mn
		requestC := make(chan *Request, rbs)
		m.requestC = requestC
		quitC := make(chan struct{})
		m.quitC = quitC
		quitErrC := make(chan error, ebs)
		exitC := make(chan struct{})
		m.exitC = exitC
		logQuitErrDoneC := make(chan struct{})
		m.logQuitErrDoneC = logQuitErrDoneC
		reqCs := make([]chan<- *Request, 0, mn)
		m.wg.Add(mn)
		for i := 0; i < mn; i++ {
			reqC := make(chan *Request, rbs)
			reqCs = append(reqCs, reqC)
			go m.matcherProcess(i, reqC, quitC, quitErrC, exitC)
		}
		m.wg.Add(1)
		go m.broadcastRequest(requestC, reqCs, quitC, quitErrC, exitC)
		go m.closeQuitErrorChannel(quitErrC)
		go m.logQuitError(quitErrC, logQuitErrDoneC)
	})
	if isInitialized {
		m.lock.RLock()
		defer m.lock.RUnlock()
		if matcherNumber > 0 && matcherNumber != m.matcherNumber {
			return ErrMatcherAlreadyInit
		}
	}
	return nil
}

func (m *Matcher) Exit(mode ExitMode) <-chan struct{} {
	exitDoneC := make(chan struct{})
	if m == nil {
		close(exitDoneC)
		return exitDoneC
	}
	isInitialized := true
	// It makes sure that cannot init after exit.
	m.initOnce.Do(func() {
		isInitialized = false
	})
	isExited := true
	m.exitOnce.Do(func() {
		isExited = false
		if !isInitialized {
			close(exitDoneC)
			return
		}
		go m.exitProcess(mode, exitDoneC)
	})
	if isExited {
		close(exitDoneC)
	}
	return exitDoneC
}

func (m *Matcher) GetMatcherNumber() int {
	if m == nil {
		return 0
	}
	err := m.Init(-1)
	if err != nil {
		return 0
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.matcherNumber
}

func (m *Matcher) IsValid() bool {
	if m == nil {
		return false
	}
	err := m.Init(-1)
	if err != nil {
		return false
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.requestC != nil
}

func (m *Matcher) SendRequest(req *Request) error {
	if m == nil {
		return ErrNilMatcher
	}
	if req == nil {
		return ErrNilRequest
	}
	err := m.Init(-1)
	if err != nil {
		return err
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.requestC == nil {
		return ErrMatcherShutdown
	}
	m.requestC <- req
	return nil
}

func (m *Matcher) broadcastRequest(inC <-chan *Request, outCs []chan<- *Request,
	quitC <-chan struct{}, quitErrC chan<- error, exitC <-chan struct{}) {
	if m == nil {
		panic(ErrNilMatcher)
	}
	defer m.wg.Done()
	// Delete nil out channels.
	outChannels := make([]chan<- *Request, 0, len(outCs))
	for _, outC := range outCs {
		if outC != nil {
			outChannels = append(outChannels, outC)
		}
	}
	isQuit := false
	for !isQuit {
		select {
		case <-exitC:
			return
		case <-quitC:
			isQuit = true
			rqn := len(inC)
			if rqn > 0 && quitErrC != nil {
				quitErrC <- fmt.Errorf(
					"broadcaster - remaining requests: %d", rqn)
			}
		case req, ok := <-inC:
			if ok {
				if req == nil {
					break
				}
				func() {
					defer req.Wg.Done()
					req.Wg.Add(len(outChannels))
					for _, outC := range outChannels {
						outC <- req
					}
				}()
			} else {
				isQuit = true
				for _, outC := range outChannels {
					close(outC)
				}
			}
		}
	}
}

func (m *Matcher) closeQuitErrorChannel(quitErrC chan<- error) {
	if quitErrC == nil {
		panic(errors.New("quit error channel is nil"))
	}
	defer close(quitErrC)
	if m == nil {
		panic(ErrNilMatcher)
	}
	m.wg.Wait()
}

func (m *Matcher) logQuitError(quitErrC <-chan error, doneC chan<- struct{}) {
	if doneC != nil {
		defer close(doneC)
	}
	if m == nil {
		panic(ErrNilMatcher)
	}
	for err := range quitErrC {
		log.Println(err)
	}
}

func (m *Matcher) matcherProcess(number int, reqC <-chan *Request,
	quitC <-chan struct{}, quitErrC chan<- error, exitC <-chan struct{}) {
	if m == nil {
		panic(ErrNilMatcher)
	}
	if reqC == nil {
		panic(errors.New("request channel is nil"))
	}
	defer m.wg.Done()
	isQuit := false
	for !isQuit {
		select {
		case <-exitC:
			return
		case <-quitC:
			isQuit = true
			rqn := len(reqC)
			if rqn > 0 && quitErrC != nil {
				quitErrC <- fmt.Errorf(
					"matcher No.%d - remaining requests: %d", number, rqn)
			}
		case req, ok := <-reqC:
			if ok {
				m.dispatchRequest(number, req, quitC, exitC)
			} else {
				isQuit = true
			}
		}
	}
}

func (m *Matcher) dispatchRequest(number int, req *Request,
	quitC <-chan struct{}, exitC <-chan struct{}) {
	if req == nil {
		return
	}
	defer req.Wg.Done()
	var sendErr interface{}
	defer func() {
		if sendErr != nil && req.ErrC != nil {
			req.ErrC <- fmt.Errorf("matcher No.%d - %v", number, sendErr)
		}
	}()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			sendErr = panicErr
		}
	}()
	if m == nil {
		sendErr = ErrNilMatcher
		return
	}
	if req.Data == nil {
		sendErr = ErrReqNilData
		return
	}
	if req.InC == nil {
		sendErr = ErrReqNilInC
		return
	}
	if req.OutC == nil {
		sendErr = ErrReqNilOutC
		return
	}
	var candidateBuffer *topkbuf.TopKBuffer
	if req.TopNumber > 0 {
		var err error
		candidateBuffer, err = topkbuf.NewTopKBuffer(req.TopNumber)
		if err != nil {
			sendErr = err
			return
		}
	}
	isQuit := false
	for !isQuit {
		select {
		case <-exitC:
			return
		case <-req.ExitC:
			return
		case <-quitC:
			isQuit = true
		case <-req.QuitC:
			isQuit = true
		case qv, ok := <-req.InC:
			if !ok {
				isQuit = true
				break
			}
			score := req.Data.Cosine(qv.TitleVector)
			if math.IsNaN(float64(score)) {
				break
			}
			candidate := &Candidate{
				QuestionId: qv.QuestionId,
				Similarity: score,
			}
			if candidateBuffer != nil {
				candidateBuffer.Add(candidate)
			} else {
				req.OutC <- candidate
			}
		}
	}
	if candidateBuffer != nil {
		outs, err := candidateBuffer.Flush()
		if err != nil {
			sendErr = err
			return
		}
		for _, out := range outs {
			candidate, ok := out.(*Candidate)
			if !ok {
				sendErr = container.ErrWrongType
				return
			}
			req.OutC <- candidate
		}
	}
}

func (m *Matcher) exitProcess(mode ExitMode, exitDoneC chan<- struct{}) {
	if exitDoneC != nil {
		defer close(exitDoneC)
	}
	if m == nil {
		panic(ErrNilMatcher)
	}
	if !mode.IsValid() {
		mode = GlobalSettings.DefaultExitMode
		if !mode.IsValid() {
			mode = ExitModeImmediately
		}
		log.Printf("Exit mode is unknown. Use default mode %v instead.\n", mode)
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	defer func() {
		m.logQuitErrDoneC = nil
		m.matcherNumber = 0
	}()
	defer func() {
		if m.exitC != nil {
			close(m.exitC)
			m.exitC = nil
		}
	}()
	defer func() {
		if m.quitC != nil {
			close(m.quitC)
			m.quitC = nil
		}
	}()
	defer func() {
		if m.requestC != nil {
			close(m.requestC)
			m.requestC = nil
		}
	}()
	switch mode {
	case ExitModeGracefully:
		close(m.requestC)
		m.requestC = nil
	case ExitModeImmediately:
		close(m.quitC)
		m.quitC = nil
	case ExitModeForcedly:
		close(m.exitC)
		m.exitC = nil
	default:
		log.Printf("Exit mode is unsupported. Use %v instead.\n",
			ExitModeImmediately)
		close(m.quitC)
		m.quitC = nil
	}
	m.wg.Wait()
	<-m.logQuitErrDoneC
}
