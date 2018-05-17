package main

import (
	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/tool/classifier"
)

func main() {
	defer db.CleanUpSessionPool()
	classifier.ClassifyByTag()
}
