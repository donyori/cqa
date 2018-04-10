package topkbuf

import (
	"testing"

	"github.com/donyori/cqa/common/container"
	"github.com/donyori/cqa/common/container/cmp"
)

type testElement int

func TestTopKBuffer(t *testing.T) {
	tkb, err := NewTopKBuffer(-2)
	if err != nil {
		if err == ErrNonPositiveK {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	tkb, err = NewTopKBuffer(5)
	if err != nil {
		t.Fatal(err)
	}
	inputs := []testElement{3, 0, 9, -4, 3, -5, 8}
	for i := range inputs {
		tkb.Add(&inputs[i])
		l := tkb.Len()
		k := tkb.K()
		shouldL := i + 1
		if shouldL > k {
			shouldL = k
		}
		if l == shouldL {
			t.Logf("tkb.Len() = %d", l)
		} else {
			t.Fatalf("tkb.Len(): %d != %d", l, shouldL)
		}
	}
	outputs, err := tkb.Flush()
	if err != nil {
		t.Fatal(err)
	}
	for i, x := range outputs {
		v := x.(*testElement)
		t.Logf("%d: %v", i, *v)
	}
}

func (te *testElement) Less(another cmp.Comparable) (
	res bool, err error) {
	a, ok := another.(*testElement)
	if !ok {
		return false, container.ErrWrongType
	}
	if te == nil {
		return a != nil, nil
	} else if a == nil {
		return false, nil
	} else {
		return *te < *a, nil
	}
}
