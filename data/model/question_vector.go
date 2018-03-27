package model

type QuestionVector struct {
	QuestionId  int64     `json:"question_id" bson:"_id"`
	TitleVector *Vector32 `json:"title_vector" bson:"title_vector"`
}
