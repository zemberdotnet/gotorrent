package coordinator

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"io"
	"sync"
)

var (
	_ interfaces.Piece = &abstractPiece{}
)

// Abstract piece defines a piece that works for both WebSeeding and BitTorrent downloads
type abstractPiece struct {
	piece  []byte
	hash   [][]byte
	index  int
	begin  int
	length int
	pos    int
}

// Index returns the index of the first piece
func (a *abstractPiece) Index() int {
	return a.index
}

// Length returns the number of pieces
func (a *abstractPiece) Length() int {
	return a.length
}

// Begin returns the block inside the piece to start downloading from
func (a *abstractPiece) Begin() int {
	return a.begin
}

// Validate validates pieces against a sha1 hash provided in the .torrent file
func (a *abstractPiece) Validate() bool {
	if len(a.hash) != a.length {
		fmt.Println("Piece download differed in size from hashes")
		return false
	}
	pieceLen := len(a.piece) / a.length
	for i, hash := range a.hash {
		shaHash := sha1.Sum(a.piece[i*pieceLen : i*pieceLen+pieceLen])
		if bytes.Compare(shaHash[:], hash) != 0 {
			fmt.Println("Piece did not match has")
			return false
		}

	}
	return true
}

// Write implements the io.Writer interface for abstractPiece
func (a *abstractPiece) Write(b []byte) (n int, err error) {
	if a.piece == nil {
		a.piece = make([]byte, len(b))
		n = copy(a.piece[:], b[:])
	} else {
		a.piece = append(a.piece, b...)
		n = len(b)
	}
	return
}

// Read implements the io.Reader interface for abstractPiece
func (a *abstractPiece) Read(b []byte) (n int, err error) {
	n = copy(b, a.piece[a.pos:])
	a.pos += n
	if a.pos == len(a.piece) {
		return n, io.EOF
	}
	return
}

// PieceFactory returns a function that allows a work scheduler to create a piece for a strategy
func PieceFactory(bt *bitfield.Bitfield, hashes [][]byte) func(multi bool) interfaces.Piece {
	var wlock sync.Mutex
	// multi is true for strategies that can accept multi-piece downloads
	return func(multi bool) interfaces.Piece {
		if multi {

			wlock.Lock()
			defer wlock.Unlock()

			index, length := bt.LargestGap()
			if length == 0 {
				return nil
			}
			// The max-dl size is the trade-off
			// larger size allows fewer connections and less overhead
			// however the cost of a failed download becomes higher
			if length > 32 {
				bt.SetPieceRange(index, 32) // hardcoded about 16mb
				length = 32
			} else {
				bt.SetPieceRange(index, length)
			}

			newHashes := make([][]byte, 0)
			for i := 0; i < length; i++ {
				newHashes = append(newHashes, hashes[index+i])

			}
			fmt.Println(newHashes)
			fmt.Println(hashes[index : index+length])

			return &abstractPiece{
				index:  index,
				length: length,
				begin:  0,
				hash:   newHashes,
			}

		} else {
			// single piece is a WIP
			return &abstractPiece{
				index:  0,
				length: 0,
				begin:  0,
				hash:   nil,
			}

		}
	}
}
