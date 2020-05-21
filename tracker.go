package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"io"
	"log"
	"net/http"
)

// Should there be multiple levels to this, would it make the code better, easier to read
// Types in this request
type TrackReq struct {
	Announce string
	InfoHash string
	PeerId   string
	Port     string
	Up       string
	Down     string
}

type TrackResp struct {
	Failure string `bencode:"failure reason"`
	//Warning    string `bencode:"warning message"`
	Interval int `bencode:"interval"`
	//MinInter   string `bencode:"min interval"`
	//TrackerId  string `bencode:"tracker id"`
	//Complete   int    `bencode:"complete"` //seeders
	//Incomplete int    `bencode:"incomplete"`
	Peers []Peer `bencode:"peers"`
}

// Need an array of Peers
type Peer struct {
	//	PeerId string `bencode:"peer id"`
	IP string `bencode:"ip"`
	// Going to try with just regular ints, could go wrong but w/e
	Port int `bencode:"port"`
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

// We have to come back and fix this hard coding
// Maybe needs to take in some parameters
func (t TorrentInfo) NewTracker() (tr TrackReq, err error) {
	// Function body should create a new tracking request

	hash := t.Info.hash()
	// How to eliminate this kind of thing
	urlH := urlHash(hash)
	fmt.Println(urlH)

	s := "unique string"
	peerid := urlHash(sha1.Sum([]byte(s)))

	return TrackReq{
		Announce: t.Announce,
		InfoHash: urlH,
		PeerId:   peerid,
	}, nil
}
func (t TrackReq) getReq() (resp *http.Response, err error) {
	return http.Get(t.Announce + "?info_hash=" + t.InfoHash + "&peer_id=" + t.PeerId)
}

func parseResponse(r io.ReadCloser) (t TrackResp, e error) {
	tr := TrackResp{}
	tr.Failure = "nil"
	err := bencode.Unmarshal(r, &tr)
	if err != nil {
		return tr, err
	}
	if tr.Failure != "nil" {
		return tr, errors.New(tr.Failure)
	}
	return tr, nil
}
