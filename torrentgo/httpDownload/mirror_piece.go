package httpDownload

import (
	"io"
)

// Currently unused, likely to be deleted soon as abstractPiece does the job
// for all strategies

// mirrorPiece is an implementation of interfaces.Piece
type mirrorPiece struct {
	piece  []byte
	hash   string
	index  int
	begin  int
	length int
	pos    int
}

func (m *mirrorPiece) Write(b []byte) (n int, err error) {
	if m.piece == nil {
		m.piece = make([]byte, len(b))
		n = copy(m.piece[:], b[:])
	} else { // this won't work when we goto bitTorrent, here we garunetee sequential
		m.piece = append(m.piece, b...)
	}
	return n, nil
}

func (m *mirrorPiece) Read(b []byte) (n int, err error) {
	n = copy(b, m.piece[m.pos:])
	m.pos += n
	if m.pos == len(m.piece) {
		return n, io.EOF
	}
	return
}

func (m *mirrorPiece) Offset() int {
	return m.index
}

// need impementation
func (m *mirrorPiece) Validate() bool {
	return false
}

func (m *mirrorPiece) Length() int {
	return m.length
}

func (m *mirrorPiece) Begin() int {
	return m.begin
}
