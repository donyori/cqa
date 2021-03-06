package web

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/donyori/cqa/data/model"
	"github.com/donyori/cqa/qa"
)

var (
	ErrParamTopIsNotInt error = errors.New(
		"parameter top is NOT an integer")
	ErrParamTimeLimitIsNotInt64 error = errors.New(
		"parameter tl is NOT a 64-bit integer")
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	startTime := time.Now()
	q := r.FormValue("q")
	if q == "" {
		err = Render(w, "index.tmpl", nil)
		if err != nil {
			HandleInternalServerError(w, err)
		}
		return
	}
	topStr := r.FormValue("top")
	var top int
	if topStr != "" {
		top, err = strconv.Atoi(topStr)
		if err != nil {
			HandleBadRequest(w, ErrParamTopIsNotInt)
			return
		}
	} else {
		top = -1
	}
	tlmsStr := r.FormValue("tl")
	var tlms int64
	if tlmsStr != "" {
		tlms, err = strconv.ParseInt(tlmsStr, 10, 64)
		if err != nil {
			HandleBadRequest(w, ErrParamTimeLimitIsNotInt64)
			return
		}
	} else {
		tlms = -1
	}
	var tl time.Duration
	if tlms > 0 {
		tl = time.Millisecond * time.Duration(tlms)
	} else {
		tl = time.Duration(tlms)
	}
	sqs, isTimeout, err := qa.SearchSimilarQuestions(q, top, tl)
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}
	similarQuestions := make([]*SimilarQuestionData, 0, len(sqs))
	for _, sq := range sqs {
		if sq == nil || sq.Question == nil {
			continue
		}
		var aa *model.Answer
		var ba *model.Answer
		for _, a := range sq.Question.Answers {
			if a.IsAccepted {
				aa = a
			}
			if ba == nil || a.Score > ba.Score {
				ba = a
			}
		}
		if ba != nil && (aa == ba || (aa != nil && aa.Score >= ba.Score)) {
			ba = nil
		}
		sqd := &SimilarQuestionData{
			Title:      sq.Question.Title,
			Content:    template.HTML(sq.Question.BodyHtml),
			Score:      sq.Question.Score,
			Link:       sq.Question.Link,
			Similarity: sq.Similarity,
		}
		if aa != nil {
			sqd.AcceptedAnswer = &SimilarQuestionAnswerData{
				Content: template.HTML(aa.BodyHtml),
				Score:   aa.Score,
				Link:    aa.Link,
			}
		}
		if ba != nil {
			sqd.BestAnswer = &SimilarQuestionAnswerData{
				Content: template.HTML(ba.BodyHtml),
				Score:   ba.Score,
				Link:    ba.Link,
			}
		}
		similarQuestions = append(similarQuestions, sqd)
	}
	elapsed := time.Since(startTime)
	err = Render(w, "qa_result.tmpl", &QaResultData{
		Question:         q,
		Top:              top,
		TimeLimitMs:      tlms,
		SimilarQuestions: similarQuestions,
		Elapsed:          elapsed,
		IsTimeout:        isTimeout,
	})
	if err != nil {
		HandleInternalServerError(w, err)
	}
}
