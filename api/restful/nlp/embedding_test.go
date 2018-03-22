package nlp

import (
	"testing"
)

func TestEmbedding(t *testing.T) {
	vec, err := Embedding("Who is Mark Twain?")
	if err == nil {
		t.Logf("result: %v", vec)
		t.Logf("dim: %v", len(vec))
		t.Logf("cap: %v", cap(vec))
	} else {
		t.Fatal(err)
	}
}
