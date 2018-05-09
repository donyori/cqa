package embedding

import (
	"time"

	"github.com/donyori/cqa/data/model"
)

type MetaValue struct {
	QuestionLastEmbeddingTime    time.Time `json:"question_last_embedding_time" bson:"question_last_embedding_time"`
	QuestionLastCreateOrEditDate time.Time `json:"question_last_create_or_edit_date" bson:"question_last_create_or_edit_date"`
}

type Meta struct {
	model.MetaKey `json:",inline" bson:",inline" cqadm:",inline"`
	Value         *MetaValue `json:"value" bson:"value"`
}

const MetaKey string = "embedding"

func NewMetaValue() *MetaValue {
	return new(MetaValue)
}

func NewMeta() *Meta {
	m := new(Meta)
	m.Key = MetaKey
	return m
}
