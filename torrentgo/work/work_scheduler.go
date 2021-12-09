package work

/*
type WorkScheduler struct {
	state *state.TorrentState
	// finding piece quickly
	// we can use this to update count then call Fix
	workQueues map[bool]chan *piece.Piece
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
	// we reset fix the piece heap
	// worst case scenario is we have to do this for
	// each piece in which case n * log(n)
	for i := 0; i < ws.PieceHeap.Len(); i++ {
		e := ws.PieceHeap.Get(i)
		// check if the count on the piece matches the count on the global state
		if e.Count != ws.state.Counts[e.Index()] {
			e.Count = ws.state.Counts[e.Index()]
			heap.Fix(ws.PieceHeap, i)
		}
	}

	// TODO THIS IS TEMPORARY SHOULD SUPPORT MULTI AND SINGLE PIECE
	// we get the single piece
	// true - single piece
	wq := ws.workQueues[true]
	// I am worried about races here
	// if there are less pieces being processed
	// than capacity then we generate a new piece

	for canAddToWorkQueue(wq, ws.state) {
		ws.state.IncrmentInProcess()
		wq <- ws.GeneratePiece()
	}

}

func canAddToWorkQueue(wq chan *piece.Piece, state *state.TorrentState) bool {
	return state.InProcess() < cap(wq)
}

func NewWorkScheduler(state *state.TorrentState, wq map[bool]chan *piece.Piece) *WorkScheduler {
	return &WorkScheduler{
		state:      state,
		workQueues: wq,
		PieceHeap:  new(piece.PieceHeap),
	}
}

// while its not perfect abstraction, multipiece is only called on
// strategies where connections have all pieces available
func (ws *WorkScheduler) GenerateMultipiece() []*piece.Piece {
	// simple linear search on bitfield
	return nil
}

func (ws *WorkScheduler) GeneratePiece() *piece.Piece {
	if ws.PieceHeap.Len() > 0 {
		return heap.Pop(ws.PieceHeap).(*piece.Piece)
	}
	return nil
}

*/
