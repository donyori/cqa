package nlp

import (
	"testing"
)

func TestEmbedding(t *testing.T) {
	vec, err := Embedding("Who is Mark Twain?")
	if err == nil {
		t.Logf("result: %+v", *vec)
		t.Logf("dim: %v", len(vec.Data))
		t.Logf("cap: %v", cap(vec.Data))
	} else {
		t.Fatal(err)
	}
}
