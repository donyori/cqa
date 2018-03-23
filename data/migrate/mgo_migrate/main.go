package main

import (
	"runtime"

	"github.com/donyori/cqa/data/migrate"
)

func main() {
	migrate.MigrateQuestions(runtime.NumCPU())
}
