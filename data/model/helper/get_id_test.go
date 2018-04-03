package helper

import (
	"testing"

	"github.com/donyori/cqa/data/model"
)

func TestGetId(t *testing.T) {
	question := model.NewQuestion()
	question.QuestionId = 123
	pQ := &question
	ppQ := &pQ
	cases := []struct {
		Input  interface{}
		Output interface{}
		Error  error
	}{
		{question, int64(123), nil},
		{pQ, int64(123), nil},
		{ppQ, int64(123), nil},
		{nil, nil, ErrNilModel},
		{model.NewVector32(), nil, ErrUnknownModelType},
	}
	for _, c := range cases {
		out, err := GetId(c.Input)
		if out != c.Output {
			t.Fatalf("Out: %v != %v", out, c.Output)
		}
		if err != c.Error {
			t.Fatalf("Out: %v != %v", err, c.Error)
		}
	}
}