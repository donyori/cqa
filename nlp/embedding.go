package nlp

import (
	api "github.com/donyori/cqa/api/restful/nlp"
	"github.com/donyori/cqa/data/dtype"
)

func Embedding(doc string) (vector *dtype.Vector32, err error) {
	switch GlobalSettings.EmbeddingMethod {
	case NlpMethodUseRestfulApi:
		return api.Embedding(doc)
	default:
		return nil, ErrUnknownNlpMethod
	}
}
