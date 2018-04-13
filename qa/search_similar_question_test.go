package qa

import (
	"sync"
	"testing"
	"time"
)

func TestSearchSimilarQuestions(t *testing.T) {
	defer func() {
		Shutdown()
		t.Log("Shutdown successfully.")
	}()
	Init()
	t.Log("Init successfully.")
	questions := [...]string{
		"What is pointer?",
		"What's the cost of vector?",
		"What is the biggest int in C?",
	}
	outs := make([]struct {
		res []*SimilarQuestion
		err error
	}, len(questions))
	var wg sync.WaitGroup
	wg.Add(len(questions))
	startTime := time.Now()
	for i := range questions {
		go func(number int) {
			defer wg.Done()
			t.Logf("%d start.", number)
			outs[number].res, outs[number].err = SearchSimilarQuestions(
				questions[number], -1, -1)
			t.Logf("%d done.", number)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(startTime)
	t.Logf("Execution time: %v", elapsed)
	for i := range outs {
		t.Logf("%d:", i)
		t.Log("  res:")
		for _, r := range outs[i].res {
			t.Logf("    %+v", *r)
		}
		t.Log("  err =", outs[i].err)
	}
}
