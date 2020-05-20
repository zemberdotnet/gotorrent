package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
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
	Failure    string `bencode:"failure reason"`
	Warning    string `bencode:"warning message"`
	Interval   string `bencode:"interval"`
	MinInter   string `bencode:"min interval"`
	TrackerId  string `bencode:"tracker id"`
	Complete   int    `bencode:"complete"` //seeders
	Incomplete int    `bencode:"incomplete"`
	Peers      []Peer `bencode:"peers"`
}

// Need an array of Peers
type Peer struct {
	PeerId string `bencode:"peer id"`
	IP     string `bencode:"ip"`
	Port   string `bencode:"port"`
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

/*
WIP PARSE THE RESPONSE
func parseReponse(resp *http.Response) TrackResp {
	resp.Body()
*/
