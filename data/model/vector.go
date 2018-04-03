package model

import (
	stdmath "math"

	"github.com/donyori/cqa/math.32"
)

type Vector32 struct {
	Data   []float32 `json:"data" bson:"data"`
	L2Norm float32   `json:"l2_norm" bson:"l2_norm"`
}

func NewVector32() *Vector32 {
	return new(Vector32)
}

func (v *Vector32) DotProduct(vector *Vector32) float32 {
	if v == nil || vector == nil || len(v.Data) == 0 || len(vector.Data) == 0 {
		return 0.0
	}
	return math.DotProduct(v.Data, vector.Data)
}

func (v *Vector32) Cosine(vector *Vector32) float32 {
	if v == nil || vector == nil || len(v.Data) == 0 || len(vector.Data) == 0 {
		return float32(stdmath.NaN())
	}
	dp := math.DotProduct(v.Data, vector.Data)
	norm1 := v.L2Norm
	norm2 := vector.L2Norm
	if norm1 == 0.0 {
		norm1 = math.L2Norm(v.Data)
	}
	if norm2 == 0.0 {
		norm2 = math.L2Norm(vector.Data)
	}
	return dp / (norm1 * norm2)
}
