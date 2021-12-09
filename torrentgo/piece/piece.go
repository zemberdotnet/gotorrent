package piece

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"io"
)

type Piece interface {
	// the index of our piece
	Index() int
	// the length in bytes of our piece
	Length() int
	// where the piece begins within the larger piece
	// i.e. the offset
	Begin() int
	// The number of bytes we have downloaded
	Verify() error
	io.Writer
	io.Reader
}

var (
	_ Piece = NewPiece(10, 0, 0, []byte{})
)

type TorrPiece struct {
	Piece      []byte
	index      int
	begin      int
	length     int
	Downloaded int
	Requested  int
	// for comparison
	Count int
	// for reading
	pos  int
	hash []byte
}

func NewPiece(size, index, begin int, hash []byte) *TorrPiece {
	return &TorrPiece{
		Piece:  make([]byte, size),
		length: size,
		index:  index,
		begin:  begin,
		hash:   hash,
	}
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

func (p *TorrPiece) Verify() error {
	hash := sha1.Sum(p.Piece)
	if !bytes.Equal(hash[:], p.hash[:]) {
		return errors.New("hash does not match")
	}

	return nil
}

/**********************************
if we want this to work when failing begin will need to move as it is written
**********************************/

// WriterAt implments the io.WriterAt interface for Piece
// NOTE: We will extend the underlying piece to be able to handle offset of any size
// This can be dangerous if the offset is incorrect
// We do not count the empty bytes written to the piece
func (p *TorrPiece) WriteAt(b []byte, off int) (n int, err error) {
	if off < 0 {
		return 0, errors.New("offset less than zero")
	}

	if p.Piece == nil {
		p.Piece = make([]byte, len(b)+off)
		n = copy(p.Piece[off:], b[:])
	} else if len(p.Piece) < off {
		diff := off - len(p.Piece)
		// here we extend the array to handle the new bytes
		p.Piece = append(p.Piece, make([]byte, diff)...)
		n = copy(p.Piece[off:], b[:])
	} else {
		n = copy(p.Piece[off:], b[:])
	}
	return n, nil
}

// Write implments the io.Writer interface for Piece
func (p *TorrPiece) Write(b []byte) (n int, err error) {
	if p.Piece == nil {
		p.Piece = make([]byte, len(b))
		n = copy(p.Piece[:], b[:])
	} else {
		p.Piece = append(p.Piece, b...)
		n = len(b)
	}
	return
}

// Read implements the io.Reader interface for Piece
func (p *TorrPiece) Read(b []byte) (n int, err error) {
	n = copy(b, p.Piece[p.pos:])
	p.pos += n
	if p.pos == len(p.Piece) {
		return n, io.EOF
	}
	return
}
