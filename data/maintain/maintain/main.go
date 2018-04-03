package main

import (
	"fmt"
	"os"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/maintain"
)

const (
	ActionEnsureIndexes   string = "ensure_indexes"
	ActionEnsureDataTypes string = "ensure_data_types"
)

func main() {
	defer db.CleanUpSessionPool()
	var action string
	actionUsage := fmt.Sprintf("Action should be one of %q or %q.",
		ActionEnsureIndexes, ActionEnsureDataTypes)
	if len(os.Args) >= 2 {
		action = os.Args[1]
	}
	switch action {
	case ActionEnsureIndexes:
		maintain.EnsureIndexes()
	case ActionEnsureDataTypes:
		maintain.EnsureDataTypes()
	case "":
		fmt.Printf("Please input the action to do. %s\n", actionUsage)
	default:
		fmt.Printf("Unknown action %q. %s\n", action, actionUsage)
	}
}
