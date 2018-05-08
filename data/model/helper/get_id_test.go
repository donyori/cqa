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
		{question, (model.Id)(123), nil},
		{pQ, (model.Id)(123), nil},
		{ppQ, (model.Id)(123), nil},
		{nil, nil, ErrNilModel},
		{model.NewVector32(), nil, ErrNoSuchField},
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

func TestGetMongoDbIdNested(t *testing.T) {
	meta := model.NewMetaInt64()
	meta.Key = "testkey"
	out, err := GetMongoDbId(meta)
	if err != nil {
		t.Fatal(err)
	}
	if out != meta.Key {
		t.Fatalf("Wrong key: %q", out)
	}
}
