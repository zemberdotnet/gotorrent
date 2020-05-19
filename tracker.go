package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"log"
)

// Should there be multiple levels to this, would it make the code better, easier to read
// Types in this request
type TrackReq struct {
	InfoHash string
	PeerId   string
	Port     string
	Up       string
	Down     string
}

func (i InfoDict) hash() (hsum [20]byte) {
	var buf bytes.Buffer

	err := bencode.Marshal(&buf, i)
	if err != nil {
		log.Fatal("Error Marshaling Info Dictionary\n Error:", err)
	}
	return sha1.Sum(buf.Bytes())
}

func urlHash(h [20]byte) (urlEncodedHash string) {
	s := ""
	for i, _ := range h {
		u := hex.EncodeToString(h[i : i+1])
		s += "%" + u
	}
	return s
}

// Maybe needs to take in some parameters
func (t TorrentInfo) NewTracker() (err error) {
	// Function body should create a new tracking request

	hash := t.Info.hash()
	// How to eliminate this kind of thing
	urlHash := urlHash(hash)
	fmt.Println(urlHash)
	return nil
}
