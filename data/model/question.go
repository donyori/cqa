package model

type Question struct {
	QuestionId       int64      `json:"question_id" bson:"_id"`
	Title            string     `json:"title" bson:"title"`
	BodyHTML         string     `json:"body" bson:"body"`
	BodyMarkdown     string     `json:"body_markdown" bson:"body_markdown"`
	IsAnswered       bool       `json:"is_answered" bson:"is_answered"`
	AcceptedAnswerId int64      `json:"accepted_answer_id" bson:"accepted_answer_id"`
	Answers          []*Answer  `json:"answers" bson:"answers"`
	Comments         []*Comment `json:"comments" bson:"comments"`
	Tags             []string   `json:"tags" bson:"tags"`
	Link             string     `json:"link" bson:"link"`
	Score            int32      `json:"score" bson:"score"`
	ViewCount        int64      `json:"view_count" bson:"view_count"`
}
