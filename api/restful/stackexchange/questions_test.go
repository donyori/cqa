package stackexchange

import (
	"testing"
	"time"
)

func TestQuestionsSimple(t *testing.T) {
	td := time.Now()
	min, err := time.Parse("2006-01-02 15:04:05", "2018-01-01 00:00:00")
	if err != nil {
		t.Fatal(err)
	}
	res, err := Questions(GlobalSettings.StartPage, GlobalSettings.MaxPageSize,
		nil, &td, QuestionsSortCreation, OrderAsc, min, nil, "c")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("res: %+v", res)
	items := res.ExtractItems()
	t.Logf("e0: %+v", *items[0])
}
