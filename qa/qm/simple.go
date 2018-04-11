package qm

import (
	"sync"
)

var (
	defaultProvider  *Provider
	defaultMatcher   *Matcher
	defaultRequester *Requester

	simpleLock     sync.RWMutex
	simpleInitOnce sync.Once
	simpleExitOnce sync.Once
)

func Init() {
	simpleInitOnce.Do(func() {
		simpleLock.Lock()
		defer simpleLock.Unlock()
		var err error
		defaultProvider, err = NewProvider()
		if err != nil {
			panic(err)
		}
		defaultMatcher, err = NewMatcher(-1)
		if err != nil {
			panic(err)
		}
		defaultRequester, err = NewRequester(defaultProvider, defaultMatcher)
		if err != nil {
			panic(err)
		}
	})
}

func Exit(mode ExitMode) <-chan struct{} {
	var exitDoneC <-chan struct{}
	isInitialized := true
	// It makes sure that cannot init after exit.
	simpleInitOnce.Do(func() {
		isInitialized = false
	})
	simpleExitOnce.Do(func() {
		if !isInitialized {
			return
		}
		simpleLock.Lock()
		defer simpleLock.Unlock()
		defer func() {
			defaultRequester = nil
			defaultMatcher = nil
			defaultProvider = nil
		}()
		if defaultMatcher != nil {
			exitDoneC = defaultMatcher.Exit(mode)
		}
	})
	if exitDoneC == nil {
		edc := make(chan struct{})
		exitDoneC = edc
		close(edc)
	}
	return exitDoneC
}

func IsValid() bool {
	Init()
	simpleLock.RLock()
	defer simpleLock.RUnlock()
	return defaultRequester != nil
}

func GetMatcherNumber() int {
	Init()
	simpleLock.RLock()
	defer simpleLock.RUnlock()
	if defaultMatcher == nil {
		return 0
	}
	return defaultMatcher.GetMatcherNumber()
}

func Match(question string, topNumber int) (respC <-chan *Response, err error) {
	Init()
	simpleLock.RLock()
	defer simpleLock.RUnlock()
	if defaultRequester == nil {
		return nil, ErrMatcherShutdown
	}
	return defaultRequester.Match(question, topNumber)
}
