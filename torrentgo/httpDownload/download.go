package httpDownload

import (
	"fmt"

	"github.com/zemberdotnet/gotorrent/interfaces"
)

var PieceLength int

type mirrorDownload struct {
	fileLength         int
	pieceLength        int
	pieceChannel       chan interface{} // interface channel
	recieveWorkChannel chan interfaces.Work
	returnWorkChannel  chan interfaces.Work
}

// NewMirrorDownload reaturns a pointer to a mirrorDownload which is a implementation
// of interfaces.Strategy
func NewMirrorDownload(fileLength, pieceLength int) *mirrorDownload {
	PieceLength = pieceLength
	recv := make(chan interfaces.Work, 30)
	rtrn := make(chan interfaces.Work)
	return &mirrorDownload{
		fileLength:         fileLength,
		pieceLength:        pieceLength,
		recieveWorkChannel: recv,
		returnWorkChannel:  rtrn,
	}
}

// These methods provide cross-channel communication
func (d *mirrorDownload) SetPieceChannel(c chan interface{}) {
	d.pieceChannel = c
}

func (d *mirrorDownload) SetRecieveChannel(c chan interfaces.Work) {
	d.recieveWorkChannel = c
}

func (d *mirrorDownload) SetReturnChannel(c chan interfaces.Work) {
	d.returnWorkChannel = c
}

func (d *mirrorDownload) RecieveWorkChannel() chan interfaces.Work {
	return d.recieveWorkChannel
}

func (d *mirrorDownload) ReturnWorkChannel() chan interfaces.Work {
	return d.returnWorkChannel
}

func (d *mirrorDownload) PieceChannel() chan interface{} {
	return d.pieceChannel
}

// These define properties to help generate work and pieces for the strategy
func (d *mirrorDownload) Multipiece() bool {
	return true
}

func (d *mirrorDownload) URL() bool {
	return true
}

// Download is the main method for the httpDownload strategy
// it recieves work, attempts to download it, and sends the piece
// to the filebuilder
func (d mirrorDownload) Download() {

	work := <-d.recieveWorkChannel
	var conn interfaces.Connection

	piece := work.GetPiece()
	if piece == nil {
		d.returnWorkChannel <- work
		return
	}

	for conn == nil {
		conn = work.GetConnection()
	}

	b, err := conn.AttemptDownloadPiece(piece)
	if err != nil {
		fmt.Printf("Error while downloading piece:%v\n", err)
		d.pieceChannel <- piece
		d.returnWorkChannel <- work
		return
	}
	_, err = piece.Write(b)
	if err != nil {
		fmt.Printf("Failed writing bytes to piece")
		d.pieceChannel <- piece
		d.returnWorkChannel <- work
		return
	}

	d.pieceChannel <- piece
	d.returnWorkChannel <- work
}
