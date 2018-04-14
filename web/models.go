package web

import (
	"html/template"
	"time"
)

type ErrorData struct {
	StatusCode int
	Msg        string
}

type SimilarQuestionAnswerData struct {
	Content template.HTML
	Score   int32
	Link    string
}

type SimilarQuestionData struct {
	Title          string
	Content        template.HTML
	Score          int32
	Link           string
	Similarity     float32
	AcceptedAnswer *SimilarQuestionAnswerData
	BestAnswer     *SimilarQuestionAnswerData
}

type QaResultData struct {
	Question         string
	SimilarQuestions []*SimilarQuestionData
	Elapsed          time.Duration
}
