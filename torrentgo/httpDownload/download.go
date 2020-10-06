package httpDownload

import (
	"fmt"
	"github.com/zemberdotnet/gotorrent/interfaces"
)

// Download is under heavy construction. I do not plan to retain this file as it is, but
// it served as a Proof of Concept. It succesfully downloaded files under good conditions.
// Now I plan to break this file into other packages/parts and make everything more elegant

type mirrorDownload struct {
	fileLength         int
	pieceLength        int
	pieceChannel       chan interface{} // interface channel
	recieveWorkChannel chan interfaces.Work
	returnWorkChannel  chan interfaces.Work
}

func NewMirrorDownload(fileLength, pieceLength int) *mirrorDownload {
	recv := make(chan interfaces.Work, 30)
	rtrn := make(chan interfaces.Work)
	return &mirrorDownload{
		fileLength:         fileLength,
		pieceLength:        pieceLength,
		recieveWorkChannel: recv,
		returnWorkChannel:  rtrn,
	}
}

// maybe piece import here
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

func (d *mirrorDownload) Multipiece() bool {
	return true
}

func (d *mirrorDownload) URL() bool {
	return true
}

// Very crude, but succesfully downloads as a Proof of Concept
func (d mirrorDownload) Download() {
	// get the next best piece to download
	// download the next best piece

	work := <-d.recieveWorkChannel
	var conn interfaces.Connection
	// You will regret this
	// TODO MUST CHANGE OR WILL FREEZE IN CERTAIN CASES!!!!
	for conn == nil {
		conn = work.GetConnection()
	}

	b, err := conn.AttemptDownloadPiece(work.GetPiece())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))

	//	fmt.Println(work.GetTask())

	//fmt.Println(work)
	d.returnWorkChannel <- work
	// Go do things
}
