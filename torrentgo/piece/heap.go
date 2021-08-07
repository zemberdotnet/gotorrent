package piece

import (
	"container/heap"
)

// fufills heap.Interface
var (
	_ heap.Interface = &PieceHeap{}
)

type PieceHeap []*TorrPiece

func (ph PieceHeap) Len() int {
	return len(ph)
}

func (ph PieceHeap) Get(i int) *TorrPiece {
	return ph[i]
}

func (ph PieceHeap) Less(i, j int) bool {
	if ph[i].Count == 0 {
		return false
	}
	if ph[j].Count == 0 {
		return true
	}

	return ph[i].Count < ph[j].Count
}

func (ph PieceHeap) Swap(i, j int) {
	ph[i], ph[j] = ph[j], ph[i]
}

func (ph *PieceHeap) Push(x interface{}) {
	*ph = append(*ph, x.(*TorrPiece))
}

func (ph *PieceHeap) Pop() interface{} {
	old := *ph
	n := len(old)
	// take last element
	x := old[n-1]
	*ph = old[0 : n-1]
	return x
}
