package piece

import (
	"io"
)

type Piece interface {
	Index() int
	Length() int
	Begin() int
	io.Writer
	io.Reader
}

var (
	_ Piece = NewPiece(10)
)

type TorrPiece struct {
	piece  []byte
	index  int
	begin  int
	length int
	// for comparison
	Count int
	// for reading
	pos int
}

func NewPiece(size int) *TorrPiece {
	return &TorrPiece{}
}

func (p *TorrPiece) Index() int {
	return p.index
}

func (p *TorrPiece) Begin() int {
	return p.begin
}

func (p *TorrPiece) Length() int {
	return p.length
}

/**********************************
if we want this to work when failing begin will need to move as it is written
**********************************/

// Write implments the io.Writer interface for Piece
func (p *TorrPiece) Write(b []byte) (n int, err error) {
	if p.piece == nil {
		p.piece = make([]byte, len(b))
		n = copy(p.piece[:], b[:])
	} else {
		p.piece = append(p.piece, b...)
		n = len(b)
	}
	return
}

// Read implements the io.Reader interface for Piece
func (p *TorrPiece) Read(b []byte) (n int, err error) {
	n = copy(b, p.piece[p.pos:])
	p.pos += n
	if p.pos == len(p.piece) {
		return n, io.EOF
	}
	return
}
