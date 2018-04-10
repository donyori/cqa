package model

type Answer struct {
	AnswerId     Id         `json:"answer_id" bson:"_id" cqadm:"id"`
	BodyHTML     string     `json:"body" bson:"body"`
	BodyMarkdown string     `json:"body_markdown" bson:"body_markdown"`
	IsAccepted   bool       `json:"is_accepted" bson:"is_accepted"`
	QuestionId   Id         `json:"question_id" bson:"question_id"`
	Comments     []*Comment `json:"comments" bson:"comments"`
	Tags         []string   `json:"tags" bson:"tags"`
	Link         string     `json:"link" bson:"link"`
	Score        int32      `json:"score" bson:"score"`
}

func NewAnswer() *Answer {
	return new(Answer)
}
