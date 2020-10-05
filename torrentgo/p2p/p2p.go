package p2p

import (
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/piece"
)

// 16kb
const BlockSize = 16384

var (
	_ interfaces.Strategy = &TorrentDL{}
)

type TorrentDL struct {
	fileLength      int
	pieceLength     int
	pieceChan       chan piece.Piece
	recieveWorkChan chan interfaces.Work
	returnWorkChan  chan interfaces.Work
}

func NewBitTorrentDownload(fileLength int) *TorrentDL {
	recv := make(chan interfaces.Work, 30)
	rtrn := make(chan interfaces.Work)
	return &TorrentDL{
		fileLength:      fileLength,
		recieveWorkChan: recv,
		returnWorkChan:  rtrn,
	}
}

func (d *TorrentDL) Download() {
	work := <-d.recieveWorkChannel
	// once again this is bad
	var conn interfaces.Connection
	for conn == nil {
		conn = work.GetConnection()
	}
	conn.Dial()
	conn.AttemptDownloadPiece(work.GetTask())
}

func (d *TorrentDL) SetReturnChannel(w chan interfaces.Work) {
	d.returnWorkChan = w
}

func (d *TorrentDL) RecieveWorkChannel() chan interfaces.Work {
	return d.recieveWorkChan
}

func (d *TorrentDL) ReturnWorkChannel() chan interfaces.Work {
	return d.returnWorkChan
}

func (d *TorrentDL) PieceChannel() chan piece.Piece {
	return d.pieceChan
}

func (d *TorrentDL) Multipiece() bool {
	return false
}

func (d *TorrentDL) URL() bool {
	return false
}
