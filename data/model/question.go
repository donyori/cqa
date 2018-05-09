package model

import (
	"time"
)

type Question struct {
	QuestionId           Id         `json:"question_id" bson:"_id" cqadm:"id"`
	Title                string     `json:"title" bson:"title"`
	BodyHtml             string     `json:"body" bson:"body"`
	BodyMarkdown         string     `json:"body_markdown" bson:"body_markdown"`
	IsAnswered           bool       `json:"is_answered" bson:"is_answered"`
	AcceptedAnswerId     Id         `json:"accepted_answer_id" bson:"accepted_answer_id"`
	Answers              []*Answer  `json:"answers" bson:"answers"`
	Comments             []*Comment `json:"comments" bson:"comments"`
	Tags                 []string   `json:"tags" bson:"tags"`
	Link                 string     `json:"link" bson:"link"`
	Score                int32      `json:"score" bson:"score"`
	ViewCount            int64      `json:"view_count" bson:"view_count"`
	CreationDate         *time.Time `json:"creation_date" bson:"creation_date"`
	LastActivityDate     *time.Time `json:"last_activity_date" bson:"last_activity_date"`
	LastEditDate         *time.Time `json:"last_edit_date" bson:"last_edit_date"`
	LastCreateOrEditDate *time.Time `json:"last_create_or_edit_date" bson:"last_create_or_edit_date"`
}

func NewQuestion() *Question {
	return new(Question)
}
