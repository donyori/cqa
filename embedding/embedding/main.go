package main

import (
	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/embedding"
)

func main() {
	defer db.CleanUpSessionPool()
	embedding.QuestionEmbedding()
}
