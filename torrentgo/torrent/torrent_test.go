package torrent

import (
	"fmt"
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
}

func TestInfoHash(t *testing.T) {
	/*
		m := &MetaInfo{
			Info: InfoDict{
				Length: 10
				Name:
	*/

	m := &MetaInfo{}
	fmt.Println(m)
	if m.hash() != [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Hash returned non-zero hash on empty InfoDict: %v", m.hash())
	}
}
