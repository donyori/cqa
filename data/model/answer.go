package model

import (
	"time"
)

type Answer struct {
	AnswerId             Id         `json:"answer_id" bson:"_id" cqadm:"id"`
	BodyHtml             string     `json:"body" bson:"body"`
	BodyMarkdown         string     `json:"body_markdown" bson:"body_markdown"`
	IsAccepted           bool       `json:"is_accepted" bson:"is_accepted"`
	QuestionId           Id         `json:"question_id" bson:"question_id"`
	Comments             []*Comment `json:"comments" bson:"comments"`
	Tags                 []string   `json:"tags" bson:"tags"`
	Link                 string     `json:"link" bson:"link"`
	Score                int32      `json:"score" bson:"score"`
	CreationDate         *time.Time `json:"creation_date" bson:"creation_date"`
	LastActivityDate     *time.Time `json:"last_activity_date" bson:"last_activity_date"`
	LastEditDate         *time.Time `json:"last_edit_date" bson:"last_edit_date"`
	LastCreateOrEditDate *time.Time `json:"last_create_or_edit_date" bson:"last_create_or_edit_date"`
}

func NewAnswer() *Answer {
	return new(Answer)
}
