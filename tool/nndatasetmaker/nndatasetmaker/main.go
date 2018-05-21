package main

import (
	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/tool/nndatasetmaker"
)

func main() {
	defer db.CleanUpSessionPool()
	nndatasetmaker.SelectQuestions()
}
