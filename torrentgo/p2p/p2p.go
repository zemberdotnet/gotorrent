package p2p

import (
	"github.com/zemberdotnet/gotorrent/connection"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/piece"
	"github.com/zemberdotnet/gotorrent/state"
)

type P2P struct {
	state     *state.TorrentState
	connPool  connection.ConnectionPool
	pieceChan chan *piece.TorrPiece
}

func (p *P2P) Start() {

}

func (p *P2P) URL() bool {
	return false
}

func (p *P2P) Multipiece() bool {
	return false
}

func NewP2PFactory(cp connection.ConnectionPool, ch chan *piece.TorrPiece) func() interfaces.Strategy {
	return func() interfaces.Strategy {
		return &P2P{
			connPool:  cp,
			pieceChan: ch,
		}
	}
}
