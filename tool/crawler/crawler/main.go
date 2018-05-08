package main

import (
	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/tool/crawler"
)

func main() {
	defer db.CleanUpSessionPool()
	crawler.CrawlQuestions()
}
