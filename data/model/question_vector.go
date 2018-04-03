package model

type QuestionVector struct {
	QuestionId  int64     `json:"question_id" bson:"_id" cqadm:"id"`
	TitleVector *Vector32 `json:"title_vector" bson:"title_vector"`
}

func NewQuestionVector() *QuestionVector {
	return new(QuestionVector)
}