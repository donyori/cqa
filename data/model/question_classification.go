package model

type QuestionClassification struct {
	QuestionId          Id       `json:"question_id" bson:"_id" cqadm:"id"`
	ClassificationByTag []string `json:"classification_by_tag" bson:"classification_by_tag"`
	ClassificationByNn  []string `json:"classification_by_nn" bson:"classification_by_nn"`
}

func NewQuestionClassification() *QuestionClassification {
	return new(QuestionClassification)
}
