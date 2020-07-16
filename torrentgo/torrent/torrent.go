package torrent

import (
	bencode "github.com/jackpal/bencode-go"
	"io"

	"bytes"
	"crypto/sha1"
	"fmt"
)

type MetaInfo struct {
	Announce  string   `bencode:"announce"`
	Comment   string   `bencode:"comment"`
	Creation  int      `bencode:"creation date"`
	CreatedBy string   `bencode:"created by"`
	Encoding  string   `bencode:"encoding"`
	Info      InfoDict `bencode:"info"`
	InfoHash  [20]byte
}

type InfoDict struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

// Unmarshal takes in bencoded torrent informationa and returns
// a struct of Go types representing the bencoded information
func Unmarshal(f io.Reader) (Meta *MetaInfo, e error) {
	m := MetaInfo{}
	err := bencode.Unmarshal(f, &m)
	if err != nil {
		// TODO Error Handling
		return &m, err
	}
	m.InfoHash = m.hash()
	return &m, err
}

// hash re-Marshals the InfoDict of a torrent's MetaInfo and returns the
// sha-1 hash
func (m *MetaInfo) hash() [20]byte {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, m.Info)
	if err != nil {
		// TODO Handle error
	}
	fmt.Println(string(buf.Bytes()))
	return sha1.Sum(buf.Bytes())
}

func (m *MetaInfo) GetAnnounce() (url string) {
	return m.Announce
}

func (m *MetaInfo) GetHash() (hash [20]byte) {
	return m.InfoHash
}

// Define Get Methods for InfoDict??
