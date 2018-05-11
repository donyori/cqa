package mongodb

import (
	"testing"

	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/model"
)

func TestIsExistedById(t *testing.T) {
	cases := []struct {
		Id  model.Id
		Res bool
	}{
		{64689, true},
		{2, false},
	}
	sess, err := NewSession(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()
	accessor, err := NewAccessor(sess)
	if err != nil {
		t.Error(err)
		return
	}
	for i, c := range cases {
		res, err := accessor.IsExistedById(dbid.QuestionCollection, c.Id)
		if err != nil {
			t.Error("case", i, err)
		}
		if res == c.Res {
			t.Log("case", i, "pass")
		} else {
			t.Error("case", i, "fail")
		}
	}
}

func TestFetchAllByIds(t *testing.T) {
	ids := []model.Id{330, 1982, 263, 2}
	sess, err := NewSession(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()
	accessor, err := NewAccessor(sess)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := accessor.FetchAllByIds(dbid.QuestionCollection, ids, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("res =", res)
	t.Log("len(res) =", len(res.([]*model.Question)))
	for i, r := range res.([]*model.Question) {
		t.Logf("  *res[%d] = %+v", i, *r)
	}
}
