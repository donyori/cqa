package nlp

import (
	"testing"
)

func TestEmbedding(t *testing.T) {
	vec, err := Embedding("Who is Mark Twain?")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("result: %+v", *vec)
	t.Logf("dim: %d", len(vec.Data))
	t.Logf("cap: %d", cap(vec.Data))
}

func TestEmbeddingWithTokens(t *testing.T) {
	vec, tvs, err := EmbeddingWithTokens("Who is Mark Twain?")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("result: %+v", *vec)
	t.Logf("dim: %d", len(vec.Data))
	t.Logf("cap: %d", cap(vec.Data))
	t.Logf("len(tvs): %d", len(tvs))
	if len(tvs) > 0 {
		t.Logf("tvs[0]: %+v", *tvs[0])
	}
}
