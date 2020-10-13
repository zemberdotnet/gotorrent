package coordinator

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"math"
)

// A better implementation may be pulling apart the scheduling and
// production of work

type workScheduler interface {
	//Schedule(interfaces.Strategy)
	GetWork() interfaces.Work
	Return(interfaces.Work)
}

type BasicScheduler struct {
	workChan       chan interfaces.Work
	strategyToWork map[interfaces.Strategy]interfaces.Work
	pieceFactory   func(bool) interfaces.Piece
	connCreator    interfaces.ConnectionCreator
	bitfield       *bitfield.Bitfield
}

func NewBasicScheduler(bt *bitfield.Bitfield, cc interfaces.ConnectionCreator, hashes [][]byte) *BasicScheduler {
	c := make(chan interfaces.Work)
	s := make(map[interfaces.Strategy]interfaces.Work)

	return &BasicScheduler{
		workChan:       c,
		strategyToWork: s,
		pieceFactory:   PieceFactory(bt, hashes),
		connCreator:    cc,
		bitfield:       bt,
	}
}

func (b *BasicScheduler) AddStrategyToWork(s interfaces.Strategy, w interfaces.Work) {
	b.strategyToWork[s] = w
}

// Placeholder
func (b *BasicScheduler) Return(w interfaces.Work) {
}

// GetWork wraps around generateWork and sends the work to the channel if its not full
// If the queue is full then we throw away the piece and unset it in the bitfield
func (b *BasicScheduler) GetWork() interfaces.Work {
	work := b.generateWork()
	if work == nil {
		return nil
	}

	select {
	case work.GetStrategy().RecieveWorkChannel() <- work:
	default:
		b.bitfield.UnsetPieceRange(work.GetPiece().Index(), work.GetPiece().Length())

	}

	return work
}

// generateWork choses the strategy to generate work for
func (b *BasicScheduler) generateWork() interfaces.Work {
	var minStrat interfaces.Strategy
	var minQueueSize int = math.MaxInt64

	// An alternative for the future is doling out the work here
	// This would make sending on the channel easier

	for k, _ := range b.strategyToWork {
		if minStrat == nil {
			minStrat = k
			minQueueSize = len(k.RecieveWorkChannel())
		}

		if len(k.RecieveWorkChannel()) < minQueueSize {
			minQueueSize = len(k.RecieveWorkChannel())
			minStrat = k
		}
	}

	return b.NewWorkFromStrategy(minStrat)
}

func (b *BasicScheduler) NewWorkFromStrategy(s interfaces.Strategy) interfaces.Work {
	if s.Multipiece() {
		piece := b.pieceFactory(true)
		if piece == nil {
			return nil
		}

		return &AbstractWork{
			strategy: s,
			piece:    piece,
			conn:     b.connCreator,
		}
	} else {
		return &AbstractWork{
			strategy: s,
			piece:    b.pieceFactory(false),
			conn:     b.connCreator,
		}
	}
}
