package tracker

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"github.com/zemberdotnet/gotorrent/torrent"
	"io"
	"log"
	"net/http"
)

type Request struct {
	Announce   string
	InfoHash   [20]byte
	PeerID     [20]byte
	Port       int
	Uploaded   string
	Downloaded string
	Left       string
	Compact    string
}

type Response struct {
	Failure  string `bencode:"failure"`
	Interval int    `bencode:"interval"`
	//TrackerID	string
	//Complete	int
	//Incomplete	int
	Peers []Peer `bencode:"peers"`
}

type Peer struct {
	IP   string `bencode:"ip"`
	Port int    `bencode:"port"`
}

func Read(r io.ReadCloser) (trackResp *Response) {
	t := Response{}
	err := bencode.Unmarshal(r, &t)
	if err != nil {
		log.Fatalf("Error unmarshaling tracker reponse: %v", err)
	}
	return &t
}

func (t *Request) serialize() string {
	an := t.Announce
	ih := urlEncode(t.InfoHash)
	id := urlEncode(t.PeerID)
	return an + "?info_hash=" + ih + "&peer_id=" + id
}

func urlEncode(h [20]byte) (escapedString string) {
	s := ""
	for i := range h {
		u := hex.EncodeToString(h[i : i+1])
		s += "%" + u
	}
	return s
}

func Get(t *Request) (resp *http.Response, err error) {
	url := t.serialize()
	return http.Get(url)
}

func NewTracker(t *torrent.Torrent) (trackReq *Request) {

	infoHash := hash(t)
	peerid := sha1.Sum(([]byte("TATWD-----378894732")))
	tracker := Request{
		Announce: t.Announce,
		InfoHash: infoHash,
		PeerID:   peerid,
	}
	return &tracker
}

func hash(t *torrent.Torrent) [20]byte {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, t.InfoDict)
	if err != nil {
		fmt.Println(err)
	}
	return sha1.Sum(buf.Bytes())
}
