package bitfield

import (
	"fmt"
	"math"
)

type Bitfield struct {
	Bitfield        []byte // used across the dowload
	PiecesAvailable []byte // used in rarest-first
	counts          []byte // used in rarest-first
	Pieces          int
}

// NewBitfield creates a new bitfield of size n/8 to represent all the pieces
func NewBitfield(pieces int) *Bitfield {

	b := make([]byte, int(math.Ceil(float64(pieces)/8)))

	return &Bitfield{
		Bitfield: b,
		Pieces:   pieces,
	}

}

// HasPiece tells if a bitfield has a particular index set
func (bf *Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	return bf.Bitfield[byteIndex]>>(7-offset)&1 != 0
}

// SetPiece sets a bit in the bitfield
func (bf *Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	bf.Bitfield[byteIndex] |= 1 << (7 - offset)
}

// SetPieceRange sets pieces from index to index+length
func (bf *Bitfield) SetPieceRange(index, length int) {
	for i := index; i <= index+length-1; i++ {
		byteIndex := i / 8
		offset := i % 8
		bf.Bitfield[byteIndex] |= 1 << (7 - offset)
	}
}

// UnsetPieceRange unsets pieces from index to index+length
func (bf *Bitfield) UnsetPieceRange(index, length int) {
	for i := index; i <= index+length-1; i++ {
		byteIndex := i / 8
		offset := i % 8
		bf.Bitfield[byteIndex] &= 0 << (7 - offset)
	}
}

// NextRarestPiece will be the main way of identifying what piece to download
// in a BitTorrent strategy
func (b *Bitfield) NextRarestPiece() (index int) {
	return 0
}

// LargestGap returns the next largestGap in a bitfield
// Used by the httpDownload strategy
func (b *Bitfield) LargestGap() (index int, length int) {
	// A linear search is reasonably fast because we search 1/8 values
	// Could become more granular in the future, but I don't think it's
	// very necessary
	largestGap := 0
	index = 0
	currentGap := 0
	for i, e := range b.Bitfield {
		if e == 0 {
			currentGap += 1
		} else {
			if currentGap > largestGap {
				largestGap = currentGap
				index = i - currentGap
				currentGap = 0
			}
		}
	}

	// one last check for empty arrays
	if currentGap > largestGap {
		largestGap = currentGap
		index = len(b.Bitfield) - currentGap
	}

	// handling edge case around the last pieces
	if index*8 > b.Pieces {
		index = 0
		largestGap = 0
	} else if index*8+largestGap*8 > b.Pieces {
		largestGap = b.Pieces - index*8
		return index * 8, largestGap
	}

	return index * 8, largestGap * 8
}
