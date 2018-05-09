package model

type TokenVector struct {
	Text   string    `json:"text" bson:"text"`
	Vector *Vector32 `json:"vector" bson:"vector"`
}

type QuestionVector struct {
	QuestionId        Id             `json:"question_id" bson:"_id" cqadm:"id"`
	TitleVector       *Vector32      `json:"title_vector" bson:"title_vector"`
	TitleTokenVectors []*TokenVector `json:"title_token_vectors" bson:"title_token_vectors"`
}

func NewTokenVector() *TokenVector {
	return new(TokenVector)
}

func NewQuestionVector() *QuestionVector {
	return new(QuestionVector)
}
