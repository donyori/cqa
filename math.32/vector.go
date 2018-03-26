package math

import (
	stdmath "math"
)

func DotProduct(v1, v2 []float32) float32 {
	lenMin := len(v1)
	if lenMin > len(v2) {
		lenMin = len(v2)
	}
	var dp float32 = 0.0
	for i := 0; i < lenMin; i++ {
		dp += v1[i] * v2[i]
	}
	return dp
}

func L2Norm(v []float32) float32 {
	if len(v) == 0 {
		return 0.0
	}
	dp := DotProduct(v, v)
	return float32(stdmath.Sqrt(float64(dp)))
}

func Cosine(v1, v2 []float32) float32 {
	dp := DotProduct(v1, v2)
	v1Norm := L2Norm(v1)
	v2Norm := L2Norm(v2)
	return dp / (v1Norm * v2Norm)
}
