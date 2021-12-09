package work

import "github.com/zemberdotnet/gotorrent/piece"

type Work interface {
}

type SinglePieceWork struct {
	piece *piece.TorrPiece
}
