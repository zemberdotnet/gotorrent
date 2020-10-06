package bitfield

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestBitfield(t *testing.T) {
	b := NewBitfield(100)
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
			fmt.Println(i)
			t.Errorf("Piece failed to set or was overwritten")

		}
	}
}

func TestLargestGap(t *testing.T) {
	b := NewBitfield(29000)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for i := 0; i < 1000; i++ {
		randInt := r1.Intn(29000)
		if b.HasPiece(randInt) == false {
			b.SetPiece(randInt)
		}
	}

	fmt.Println(b.Bitfield)
	fmt.Println(b.LargestGap())
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
