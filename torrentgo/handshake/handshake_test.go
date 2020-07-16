package handshake

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestSerialize(t *testing.T) {
	h := &Handshake{
		InfoHash: [20]byte{65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65},
		PeerID:   [20]byte{90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90},
		Pstr:     []byte("BitTorrent protocol"),
	}

	expected := [68]byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90, 90}
	hshake := h.serialize()
	// need to read cmp doc to get this right
	if cmp.Equal(hshake, expected[:]) {
		fmt.Println("Problem")
	}

	// remove later
	for x := range hshake {
		if hshake[x] != expected[x] {
			fmt.Println(x)
		}
	}

}
