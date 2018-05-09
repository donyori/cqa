package main

import (
	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/tool/tagcount"
)

func main() {
	defer db.CleanUpSessionPool()
	tagcount.TagCount()
}
