package httpDownload

import (
	"github.com/zemberdotnet/gotorrent/connection"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/piece"
)

type MirrorDownload struct {
	connPool  *connections.ConnectionPool
	pieceChan chan *piece.Piece
}

func NewMirrorDownloadFactory(cp *connection.ConnectionPool, ch chan *piece.Piece) func() interfaces.Strategy {
	return func() interfaces.Strategy {
		return &MirrorDownload{
			connPool:  cp,
			pieceChan: ch,
		}
	}
}

func (md *MirrorDownload) Start() {
	return
}

func (md *MirrorDownload) Multipiece() bool {
	return true
}

func (md *MirrorDownload) URL() bool {
	return true
}
