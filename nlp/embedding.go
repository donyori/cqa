package nlp

import (
	api "github.com/donyori/cqa/api/restful/nlp"
	"github.com/donyori/cqa/data/model"
)

func Embedding(doc string) (vector *model.Vector32, err error) {
	switch GlobalSettings.EmbeddingMethod {
	case NlpMethodUseRestfulApi:
		return api.Embedding(doc)
	default:
		return nil, ErrUnknownNlpMethod
	}
}
