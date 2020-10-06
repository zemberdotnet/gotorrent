package bitfield

import (
	"math"
)

// Credit github.com/veggiedefender - Adding liscence next commit

// A Bitfield represents the pieces that a peer has

type Bitfield struct {
	Bitfield        []byte
	PiecesAvailable []byte
	counts          []byte
}

func NewBitfield(pieces int) *Bitfield {

	b := make([]byte, int(math.Ceil(float64(pieces)/8)))

	return &Bitfield{
		Bitfield: b,
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

func (bf *Bitfield) SetPieces(index, length int) {
	// this could cause problems
	for i := index; i <= index+length; i++ {
		byteIndex := i / 8
		offset := i % 8
		bf.Bitfield[byteIndex] |= 1 << (7 - offset)
	}
}

func (bf *Bitfield) UnsetPieces(index, length int) {
	for i := index; i <= index+length; i++ {
		byteIndex := i / 8
		offset := i % 8
		bf.Bitfield[byteIndex] |= 0 << (7 - offset)
	}
}

// could organize pieces by number of peers that have that piece
// implementation pending
func (b *Bitfield) NextRarestPiece() (index int) {
	return 0
}

// we can consider there is only two cases, each time we must search n
// elements
// n will be pieces
// until we can find something better this is what we are doing
func (b *Bitfield) LargestGap() (x int, y int) {
	largestGap := 0
	temp := 0
	for i, e := range b.Bitfield {
		if largestGap == 0 && e == 0 {
			temp = i
			largestGap++
			continue
		}

		if e == 0 {
			continue
			//	largestGap++
		} else if e != 0 {
			if i-temp > largestGap {
				largestGap = i - temp
				x = temp
				temp = i
			}
		}
	}
	if len(b.Bitfield)-1-temp > largestGap {

		largestGap = len(b.Bitfield) - 1 - temp
		x = temp
	}

	return x, largestGap
}
