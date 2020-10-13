package torrent

import (
	"bytes"
	"crypto/sha1"
	bencode "github.com/jackpal/bencode-go"
	"io"
	"log"
)

// MetaInfo represents the bencoded information in a torrent file
type MetaInfo struct {
	Announce  string   `bencode:"announce"`
	Comment   string   `bencode:"comment"`
	Creation  int      `bencode:"creation date"`
	CreatedBy string   `bencode:"created by"`
	Encoding  string   `bencode:"encoding"`
	Info      InfoDict `bencode:"info"`
	URLList   []string `bencode:"url-list"`
	InfoHash  [20]byte
}

// InfoDict represents the info dictionary in a torrent file
type InfoDict struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

// Unmarshal takes in bencoded torrent informationa and returns
// a struct of Go types representing the bencoded information
func Unmarshal(f io.Reader) (m *MetaInfo, err error) {
	m = &MetaInfo{}
	err = bencode.Unmarshal(f, m)
	if err != nil {
		return m, err
	}
	m.InfoHash = m.hash()
	return
	// return automatically returns the named parameters m and err
}

// hash re-Marshals the InfoDict of a torrent's MetaInfo and returns the
// sha-1 hash
func (m *MetaInfo) hash() [20]byte {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, m.Info)
	if err != nil {
		// TODO: Handle Error
	}
	return sha1.Sum(buf.Bytes())
}

func (m *MetaInfo) GetAnnounce() (url string) {
	return m.Announce
}

func (m *MetaInfo) GetHash() (hash [20]byte) {
	return m.InfoHash
}

func (m *MetaInfo) Length() int {
	return m.Info.Length
}

func (m *MetaInfo) Pieces() int {
	return int(float64(m.Info.Length) / float64(m.Info.PieceLength))

}

func (m *MetaInfo) PieceLength() int {
	return m.Info.PieceLength
}

func (m *MetaInfo) PieceHashes() [][]byte {
	//pieces := m.Pieces()
	// could speed up with static allocation
	pieceHashes := make([][]byte, 0)
	buf := bytes.NewBufferString(m.Info.Pieces)

	var err error
	for err != io.EOF {
		hash := make([]byte, 20)
		n, err := buf.Read(hash)
		if err == io.EOF {
			break
		}
		if err != nil || n != 20 {
			log.Fatal(err, n)
			return nil // overly dramatic
		}
		pieceHashes = append(pieceHashes, hash)
	}

	return pieceHashes

}
