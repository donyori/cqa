package nlp

import (
	"encoding/json"

	"github.com/donyori/cqa/data/model"
)

type EmbeddingDataCell struct {
	Vector       []float32 `json:"vector"`
	VectorL2Norm float32   `json:"vector_l2_norm"`
}

type EmbeddingResponse struct {
	Data EmbeddingDataCell `json:"data"`
}

type EmbeddingWithTokensResponse struct {
	Data struct {
		EmbeddingDataCell `json:",inline"`
		TokenVectors      []*struct {
			Text              string `json:"text"`
			EmbeddingDataCell `json:",inline"`
		} `json:"token_vectors"`
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
	vector = model.NewVector32()
	vector.Data = res.Data.Vector
	vector.L2Norm = res.Data.VectorL2Norm
	return vector, nil
}

func EmbeddingWithTokens(doc string) (vector *model.Vector32,
	tokenVectors []*model.TokenVector, err error) {
	resp, err := GetNlpGetClient().R().
		SetQueryParam("q", doc).
		Get(GlobalSettings.EmbeddingWithTokensPath)
	if err != nil {
		return nil, nil, err
	}
	res := new(EmbeddingWithTokensResponse)
	err = json.Unmarshal(resp.Body(), res)
	if err != nil {
		return nil, nil, err
	}
	vector = model.NewVector32()
	vector.Data = res.Data.Vector
	vector.L2Norm = res.Data.VectorL2Norm
	tokenVectors = nil
	if res.Data.TokenVectors != nil {
		tokenVectors = make([]*model.TokenVector, 0, len(res.Data.TokenVectors))
		for _, dtv := range res.Data.TokenVectors {
			tv := model.NewTokenVector()
			tv.Text = dtv.Text
			v := model.NewVector32()
			v.Data = dtv.Vector
			v.L2Norm = dtv.VectorL2Norm
			tv.Vector = v
			tokenVectors = append(tokenVectors, tv)
		}
	}
	return vector, tokenVectors, nil
}
