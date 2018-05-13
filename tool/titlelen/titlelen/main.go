package main

import (
	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/tool/titlelen"
)

func main() {
	defer db.CleanUpSessionPool()
	titlelen.TitleLengthCount()
}
