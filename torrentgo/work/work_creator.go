package work

import (
	"errors"
	"sync"

	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/piece"
	"github.com/zemberdotnet/gotorrent/state"
)

type WorkCreator struct {
	state *state.TorrentState
	// holds the indexes of pieces that are
	// already being worked on
	inProgress []int

	lock *sync.RWMutex
}

func NewWorkCreator(state *state.TorrentState) *WorkCreator {
	return &WorkCreator{
		state:      state,
		inProgress: make([]int, state.Length),
		lock:       &sync.RWMutex{},
	}
}

func (wc *WorkCreator) WorkFromBitfield(b bitfield.Bitfield, multipiece bool) ([]*piece.TorrPiece, error) {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if multipiece {
		return nil, errors.New("Not implemented")
		//return wc.createMultiPieceWork(b)
	} else {
		return wc.createSinglePieceWork(b)
	}
}

func (wc *WorkCreator) createSinglePieceWork(b bitfield.Bitfield) ([]*piece.TorrPiece, error) {
	p := -1
	min_piece := 100000000

	for i := 0; i < len(b)*8; i++ {
		if b.HasPiece(i) && wc.state.GetCount(i) < min_piece && wc.inProgress[i] == 0 {
			p = i
			min_piece = wc.state.GetCount(i)
		}
	}

	// we didn't find a suitable piece
	if p == -1 {
		return nil, errors.New("No piece is available for given bitfield")
	}
	wc.inProgress[p] = 1

	return nil, errors.New("Not implemented")
}

func (wc *WorkCreator) createMultiPieceWork(b bitfield.Bitfield) ([]Work, error) {
	return nil, errors.New("Not implemented")
}

// ReturnFailedWork is used to return work that has failed so it may be retried
func (wc *WorkCreator) ReturnFailedWork(p *piece.TorrPiece) {
	wc.inProgress[p.Index()] = 0
}
