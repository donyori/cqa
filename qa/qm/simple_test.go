package qm

import (
	"sync"
	"testing"
	"time"
)

func TestSimpleMatch(t *testing.T) {
	defer func() {
		doneC := Exit(ExitModeGracefully)
		<-doneC
		t.Log("Exit successfully.")
	}()
	Init()
	t.Log("Init successfully.")
	questions := [...]string{
		"What is pointer?",
		"What's the cost of vector?",
		"What is the biggest int in C?",
	}
	outs := make([]struct {
		resp *Response
		err  error
	}, len(questions))
	var wg sync.WaitGroup
	wg.Add(len(questions))
	startTime := time.Now()
	for i := range questions {
		go func(number int) {
			defer wg.Done()
			t.Logf("%d start.", number)
			var respC <-chan *Response
			respC, outs[number].err = Match(
				questions[number], 5, time.Second*time.Duration(number+1))
			if respC != nil {
				outs[number].resp = <-respC
			}
			t.Logf("%d done.", number)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(startTime)
	t.Logf("Execution time: %v", elapsed)
	for i, out := range outs {
		t.Logf("%d:", i)
		t.Log("  candidates:")
		for _, candidate := range out.resp.Candidates {
			t.Logf("    %+v", *candidate)
		}
		t.Log("  errors:")
		for _, e := range out.resp.Errors {
			t.Logf("    %v", e)
		}
		t.Logf("  is time out: %v", out.resp.IsTimeout)
		t.Logf("  error: %v", out.err)
	}
}
