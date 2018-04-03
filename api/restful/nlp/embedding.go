package nlp

import (
	"encoding/json"

	"github.com/donyori/cqa/data/model"
)

type EmbeddingResponse struct {
	Data struct {
		Vector       []float32 `json:"vector"`
		VectorL2Norm float32   `json:"vector_l2_norm"`
	} `json:"data"`
}

func Embedding(doc string) (vector *model.Vector32, err error) {
	resp, err := GetNlpGetClient().R().
		SetQueryParam("q", doc).
		Get(GlobalSettings.EmbeddingPath)
	if err != nil {
		return nil, err
	}
	res := new(EmbeddingResponse)
	err = json.Unmarshal(resp.Body(), res)
	if err != nil {
		return nil, err
	}
	ret := model.NewVector32()
	ret.Data = res.Data.Vector
	ret.L2Norm = res.Data.VectorL2Norm
	return ret, nil
}
