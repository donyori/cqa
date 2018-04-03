package nlp

import (
	api "github.com/donyori/cqa/api/restful/nlp"
	"github.com/donyori/cqa/data/model"
)

func Embedding(doc string) (vector *model.Vector32, err error) {
	nlpSource := GlobalSettings.EmbeddingSource
	switch nlpSource {
	case NlpSourceRestfulApi:
		return api.Embedding(doc)
	default:
		return nil, ErrInvalidNlpSource
	}
}
