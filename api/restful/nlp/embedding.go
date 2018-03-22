package nlp

import (
	"encoding/json"
)

type EmbeddingResult struct {
	Vector []float32 `json:"data"`
}

func Embedding(doc string) (vector []float32, err error) {
	resp, err := GetNlpGetClient().R().
		SetQueryParam("q", doc).
		Get(GlobalSettings.EmbeddingPath)
	if err != nil {
		return nil, err
	}
	res := new(EmbeddingResult)
	err = json.Unmarshal(resp.Body(), res)
	if err != nil {
		return nil, err
	}
	return res.Vector, nil
}
