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

/*
func TestInfoHash(t *testing.T) {
	// Empty should hash to 0s
	m := &MetaInfo{}
	// This won't work because field names are added
	if m.hash() != [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		t.Errorf("Hash returned non-zero hash on empty InfoDict: %v", m.hash())
	}

}
*/

func TestPieces(t *testing.T) {
	testData := newTestTorrent()
	_ = testData.Pieces()

}

func TestPieceHashes(t *testing.T) {
	testData := newTestTorrent()
	hashes := testData.PieceHashes()
	numOfPieces := testData.Pieces()
	if numOfPieces != len(hashes) {
		t.Errorf("Mismatch between number of pieces and number of hashes.")
	}

	for _, hash := range hashes {
		if len(hash) != 20 {
			t.Errorf("Invalid hash size %v", hash)
		}

	}

}

func newTestTorrent() *MetaInfo {
	filepath := "/home/snow/Projects/Coding/go/gotorrent/arch.torrent"
	f, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	m, err := Unmarshal(f)
	f.Close()
	if err != nil {
		return nil
	}
	return m

}
