package coordinator

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
)

var (
	_ interfaces.Piece = abstractPiece{}
)

type abstractPiece struct {
	piece   []byte
	hash    []string
	index   int
	begin   int
	length  int
	readPos int
}

func (a abstractPiece) Index() int {
	return a.index
}

func (a abstractPiece) Length() int {
	return a.length
}

func (a abstractPiece) Begin() int {
	return a.begin
}

func (a abstractPiece) Validate() bool {
	return false
}

// naive implentation, may fail on bitTorrent download because assumes
// sequential download of pieces
func (a abstractPiece) Write(b []byte) (n int, err error) {
	if a.piece == nil {
		a.piece = make([]byte, len(b))
		n = copy(a.piece[:], b[:])
	} else {
		a.piece = append(a.piece, b...)
		n = len(b)
	}
	return

}

func (a abstractPiece) Read(b []byte) (n int, err error) {
	return 0, nil
}

func PieceFactory(bt *bitfield.Bitfield) func(multi bool) interfaces.Piece {
	return func(multi bool) interfaces.Piece {
		if multi {
			index, length := bt.LargestGap()
			return abstractPiece{
				index:  index,
				length: length,
				begin:  0,
			}

		} else {
			// need to add in handling for single pieces
			return abstractPiece{
				index:  0,
				length: 0, // make this the bittorent size
				begin:  0,
			}

		}

	}
}
