package main

import (
	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/tool/embedding"
)

func main() {
	defer db.CleanUpSessionPool()
	embedding.QuestionEmbedding()
}
