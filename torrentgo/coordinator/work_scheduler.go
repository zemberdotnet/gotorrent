package coordinator

import (
	"container/heap"
	"sync"

	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/piece"
	"github.com/zemberdotnet/gotorrent/state"
)

type WorkScheduler struct {
	state *state.TorrentState
	// finding piece quickly
	// we can use this to update count then call Fix
	workQueue chan piece.Piece
	// selecting rarest piece
	PieceHeap *piece.PieceHeap
	lock      *sync.RWMutex
}

// need to scan connection pool
//
func (ws *WorkScheduler) Schedule() {
	// first scan the counts
	// and go thru pieceheap updating and fixing
	// we have to do this to use Fix anyway
	ws.lock.Lock()
	defer ws.lock.Unlock()
	// looking at a n*log(n) here
	for i := 0; i < ws.PieceHeap.Len(); i++ {
		e := ws.PieceHeap.Get(i)
		if e.Count != ws.state.Counts[e.Index()] {
			e.Count = ws.state.Counts[e.Index()]
			heap.Fix(ws.PieceHeap, i)
		}
	}

	// I am worried about races here
	for len(ws.workQueue)+ws.state.InProcess() < cap(ws.workQueue) {
		ws.workQueue <- ws.GeneratePiece()
	}

}

func NewWorkScheduler() *WorkScheduler {
	return &WorkScheduler{
		PieceHeap: new(piece.PieceHeap),
	}
}

func (ws *WorkScheduler) SendPiece(strat interfaces.Strategy) {

}

// while its not perfect abstraction, multipiece is only called on
// strategies where connections have all pieces available
func (ws *WorkScheduler) GenerateMultipiece() []*piece.Piece {
	// simple linear search on bitfield
	return nil
}

func (ws *WorkScheduler) GeneratePiece() piece.Piece {
	if ws.PieceHeap.Len() > 0 {
		return heap.Pop(ws.PieceHeap).(piece.Piece)
	}
	return nil
}
