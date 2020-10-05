package interfaces

import (
	"github.com/zemberdotnet/gotorrent/piece"
)

type Strategy interface {
	Download()
	// These are all communication features
	SetReturnChannel(chan Work)
	RecieveWorkChannel() chan Work
	ReturnWorkChannel() chan Work
	PieceChannel() chan piece.Piece
	// These are properties on the strategy that help create work for it
	Multipiece() bool
	URL() bool
}
