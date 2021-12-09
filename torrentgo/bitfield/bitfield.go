package bitfield

import (
	"math"
)

type Bitfield []byte // used across the dowload

// NewBitfield creates a new bitfield of size n/8 to represent all the pieces
func NewBitfield(pieces int) Bitfield {
	return make([]byte, int(math.Ceil(float64(pieces)/8)))
}

func NewBitfieldFromBytes(b []byte) Bitfield {
	return Bitfield(b)
}

// HasPiece tells if a bitfield has a particular index set
func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	return bf[byteIndex]>>(7-offset)&1 != 0
}

// SetPiece sets a bit in the bitfield
// Note: we rely on good input
func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	bf[byteIndex] |= 1 << (7 - offset)
}

// SetPieceRange sets pieces from index to index+length
func (bf Bitfield) SetPieceRange(index, length int) {
	for i := index; i <= index+length-1; i++ {
		byteIndex := i / 8
		offset := i % 8
		bf[byteIndex] |= 1 << (7 - offset)
	}
}

// UnsetPieceRange unsets pieces from index to index+length
func (bf Bitfield) UnsetPieceRange(index, length int) {
	for i := index; i <= index+length-1; i++ {
		byteIndex := i / 8
		offset := i % 8
		bf[byteIndex] &= 0 << (7 - offset)
	}
}

// TODO Consider where we will do the bounds checking?
// Here or elsewhere
/*
func (b *Bitfield) ParseBitfieldIntoCounts(bt []byte) {
	recvBitfield := NewBitfieldFromBytes(bt, b.Pieces)
	for i := 0; i < b.Pieces; i++ {

		if recvBitfield.HasPiece(i) {
			b.IncreaseCount(i)
		}

	}
}

func (b *Bitfield) IncreaseCount(index int) {
	b.Counts[index]++
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
*/
