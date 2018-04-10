package underlying

import (
	"github.com/donyori/cqa/common/container/cmp"
)

type Heap []cmp.Comparable

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	res, err := h[i].Less(h[j])
	if err != nil {
		panic(err)
	}
	return res
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(cmp.Comparable))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
