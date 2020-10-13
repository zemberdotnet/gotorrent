package bitfield

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestBitfield(t *testing.T) {
	b := NewBitfield(64)
	fmt.Println(b)

}

func TestSetPiece(t *testing.T) {
	size := 100
	b := NewBitfield(size)
	for i := 0; i <= size; i++ {
		b.SetPiece(i)
	}

	for i, e := range b.Bitfield {
		if e != byte(255) && i != len(b.Bitfield)-1 {
			//fmt.Println(i)
			t.Errorf("Piece failed to set or was overwritten")

		}
	}
}

// TestLargestGapFuzzing runs random data through LargestGap
func TestLargestGapFuzzing(t *testing.T) {
	b := NewBitfield(29000)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < 1000; i++ {
		randInt := r1.Intn(29000)
		if b.HasPiece(randInt) == false {
			b.SetPiece(randInt)
		}
	}

	//fmt.Println(b.Bitfield)
	fmt.Println(b.LargestGap())
}

// TestLargestGap
func TestLargestGapFull(t *testing.T) {
	fmt.Println("====LargestGapFullTest====")
	b := NewBitfield(64)
	for i := 0; i < 64; i++ {
		b.SetPiece(i)
	}
	fmt.Println(b.Bitfield)
	index, gap := b.LargestGap()
	fmt.Println(gap)
	if gap != 0 {
		t.Errorf("Largest gap produced non-zero output on full bitfield")
	}
	if index != 0 {
		t.Errorf("Largest gap produced non-ending index")
	}

}

func TestLargestGapEmpty(t *testing.T) {
	fmt.Println("====LargestGapEmpty====")
	b := NewBitfield(64)
	index, gap := b.LargestGap()
	if index != 0 {
		t.Errorf("Invalid index from empty bitfield\nIndex:%v\n", index)
	}
	if gap != 64 {
		t.Errorf("Invalid gap from empty bitfield")
	}

}

func TestLargestGapPartial(t *testing.T) {
	fmt.Println("====LargestGapPartial====")
	b := NewBitfield(65) // + 1
	index, gap := b.LargestGap()
	if index != 0 {
		t.Errorf("Invalid index from empty bitfield\nIndex:%v\n", index)
	}
	if gap != 65 {
		t.Errorf("Invalid gap from empty bitfield %v", gap)
	}

}

func TestLargestGapExample(t *testing.T) {
	b := NewBitfield(1362)
	for {
		index, gap := b.LargestGap()
		if gap > 16 {
			b.SetPieceRange(index, 16)
		} else {
			b.SetPieceRange(index, gap)
		}
		if gap == 0 {
			break
		}

	}
	fmt.Println(b.Bitfield)

}

func TestSetPiecesRange(t *testing.T) {
	fmt.Println("====TestSetPiecesRange====")
	b := NewBitfield(64)
	b.SetPiece(43)
	index, length := b.LargestGap()
	fmt.Println(b.Bitfield)
	fmt.Println(index, length)
	b.SetPieceRange(index, length)
	//fmt.Println(b.Bitfield)

}

func BenchmarkLargestGap(b *testing.B) {

	b.StopTimer()
	bt := NewBitfield(29000)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	b.StartTimer()

	for i := 0; i < 1000; i++ {
		randInt := r1.Intn(29000)
		if bt.HasPiece(randInt) == false {
			bt.SetPiece(randInt)
		}
	}
	_, _ = bt.LargestGap()
}

func TestUnsetPieceRange(t *testing.T) {
	b := NewBitfield(64)
	index, gap := b.LargestGap()
	b.SetPieceRange(index, gap)
	b.UnsetPieceRange(index, gap)
	for _, bite := range b.Bitfield {
		if bite != 0 {
			fmt.Println(b.Bitfield)
			t.Errorf("Byte was not set to zero from unset range")
		}

	}

}
