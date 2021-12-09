package work

import "github.com/zemberdotnet/gotorrent/piece"

type WorkQueue interface {
	// A work queue must tell if it
	// able to accept arbitrary length of work.
	SinglePiece() bool
}

type SinglePieceWorkQueue chan *piece.Piece
