package stackexchange

import (
	"errors"
)

type Order int8
type QuestionsSort int8

const (
	OrderAsc Order = iota
	OrderDesc
)

const (
	QuestionsSortActivity QuestionsSort = iota
	QuestionsSortCreation
	QuestionsSortVotes
	QuestionsSortHot
	QuestionsSortWeek
	QuestionsSortMonth
)

var (
	ErrUnknownOrder         error = errors.New("order param is unknown")
	ErrUnknownQuestionsSort error = errors.New(
		"questions sort param is unknown")

	orderStrings = [...]string{
		"asc",
		"desc",
	}
	questionsSortStrings = [...]string{
		"activity",
		"creation",
		"votes",
		"hot",
		"week",
		"month",
	}
)

func (o Order) IsValid() bool {
	return o >= OrderAsc && o <= OrderDesc
}

func (o Order) String() string {
	if !o.IsValid() {
		return "unknown"
	}
	return orderStrings[o]
}

func (qs QuestionsSort) IsValid() bool {
	return qs >= QuestionsSortActivity && qs <= QuestionsSortMonth
}

func (qs QuestionsSort) String() string {
	if !qs.IsValid() {
		return "unknown"
	}
	return questionsSortStrings[qs]
}
