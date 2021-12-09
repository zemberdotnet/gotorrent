package filebuilder

import (
	"testing"

	"github.com/zemberdotnet/gotorrent/piece"
)

// implements the piece interface for testing

var (
	_ piece.Piece = TestPiece{}
)

type TestPiece struct {
	index  int
	length int
}

func (t TestPiece) Index() int {
	return t.index
}

func (t TestPiece) Length() int {
	return t.length
}

func (t TestPiece) Begin() int {
	return 0
}

func (t TestPiece) Validate() bool {
	return true
}

func (t TestPiece) Write(b []byte) (int, error) {
	return 0, nil
}

func (t TestPiece) Read(b []byte) (int, error) {
	return 0, nil
}

func TestWriteToFile(t *testing.T) {
	t.Errorf("TODO")
}
