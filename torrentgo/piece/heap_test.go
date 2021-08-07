package piece

import (
	"container/heap"
	"testing"
)

// TODO Test performance tuning when both are 0
func HeapTest(t *testing.T) {
	t.Log("run")
}

func TestHeap(t *testing.T) {
	ph := &PieceHeap{}
	heap.Init(ph)
	p1 := NewPiece(0)
	p1.Count = 0
	//_ := NewPiece(0)
	//p2.Count = 2
	p3 := NewPiece(0)
	p3.Count = 3

	heap.Push(ph, p1)
	heap.Push(ph, p1)
	heap.Push(ph, p3)

	for ph.Len() > 0 {
		t.Log(heap.Pop(ph).(*TorrPiece).Count)
	}

}
