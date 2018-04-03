package helper

import (
	"github.com/donyori/cqa/data/model"
)

type Maker func() interface{}

func AnswerMaker() interface{} {
	return model.NewAnswer()
}

func CommentMaker() interface{} {
	return model.NewComment()
}

func QuestionMaker() interface{} {
	return model.NewQuestion()
}

func QuestionVectorMaker() interface{} {
	return model.NewQuestionVector()
}

func Vector32Maker() interface{} {
	return model.NewVector32()
}
