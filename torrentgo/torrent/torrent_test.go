package torrent

import (
	"os"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	filepath := "/home/snow/Projects/Coding/go/gotorrent/arch.torrent"
	f, err := os.Open(filepath)
	if err != nil {
		t.Errorf("Failed to open file at %v", filepath)
		return
	}
	_, err = Unmarshal(f)
	f.Close()
	if err != nil {
		t.Errorf("Unmarshal failed because %v", err)
	}
	// Add assertions for m

}

func TestInfoHash(t *testing.T) {
	// Empty should hash to 0s
	m := &MetaInfo{}
	// This won't work because field names are added
	if m.hash() != [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Hash returned non-zero hash on empty InfoDict: %v", m.hash())
	}

}
