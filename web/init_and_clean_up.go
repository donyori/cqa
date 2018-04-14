package web

import (
	"sync"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/qa"
)

var (
	initOnce sync.Once
)

func Init() {
	initOnce.Do(func() {
		err := LoadTemplates()
		if err != nil {
			panic(err)
		}
		qa.Init()
	})
}

func CleanUp() {
	defer CleanUpTemplates()
	defer db.CleanUpSessionPool()
	qa.Shutdown()
}
