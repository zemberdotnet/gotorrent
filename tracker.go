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

// TrackReq represents parameters needed to send a request to the tracker
// Values are encoded as strings, but it would probably be better to move these to
// bytes in a future version.
type TrackReq struct {
	Announce string
	InfoHash string
	PeerId   string
	Port     string
	Up       string
	Down     string
}

// TrackResp represents a response from the Tracker.
// NOTE: Currently commented values are optional. Future updates should reconsile these
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

// Peer represents a single peer from the peers key in the tracker response.
// A peer consists of an IP and a Port.
// NOTE: bencode passes back the port as a int and IP as a string. In the future,
// moving to big endian (unit16) and net.IP might improve things a bit
type Peer struct {
	//	PeerId string `bencode:"peer id"`
	IP string `bencode:"ip"`
	// Going to try with just regular ints, could go wrong but w/e
	Port int `bencode:"port"`
}

// hash computes the sha1 sum of the bencode formatted info key.
// We re-marshal the infodict and use the bytes to compute a [20]byte hash sum.
func (i InfoDict) hash() (hsum [20]byte) {
	var buf bytes.Buffer

	err := bencode.Marshal(&buf, i)
	if err != nil {
		log.Fatal("Error Marshaling Info Dictionary\n Error:", err)
	}
	return sha1.Sum(buf.Bytes())
}

// urlHash takes a sha1 Sum ([20]bytes) and returns a string formatted in the
// url escape sequence
func urlHash(h [20]byte) (urlEncodedHash string) {
	s := ""
	for i := range h {
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

// getReq sends the tracker request to the tracker.
// NOTE: I would like to rename this function & possibly change how parameters are used in the query
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
