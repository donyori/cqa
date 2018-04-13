package qa

import (
	"github.com/donyori/cqa/data/model"
)

type SimilarQuestion struct {
	Question   *model.Question
	Similarity float32
}
