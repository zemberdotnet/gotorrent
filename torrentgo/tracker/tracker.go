package tracker

import (
	// "fmt"
	bencode "github.com/jackpal/bencode-go"
	"github.com/zemberdotnet/gotorrent/peer"
	"github.com/zemberdotnet/gotorrent/torrent"
	"io"

	"math/rand"
	"net/http"
	"net/url"
)

// Tracker Response represents a bencoded response from a tracker
type TrackerResponse struct {
	Failure  string `bencode:"failure"`
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
	Parsed   *[]peer.Peer
}

// GetPeers
func GetPeers(m *torrent.MetaInfo) (t *TrackerResponse, e error) {
	id := newPeerID()
	// fmt.Println(m.GetHash())
	url := newRequest(m.GetAnnounce(), m.GetHash(), id)
	return getPeers(url)
}

func getPeers(url string) (tResponse *TrackerResponse, e error) {
	resp, err := http.Get(url)
	if err != nil {
		// fmt.Println("returned before unmarshalling")
		return &TrackerResponse{}, err
		// TODO handle error
	}
	defer resp.Body.Close()
	return Read(resp.Body)
}

// Format neutral to deal with bad responses from tracker
// Maybe shouldn't have this method public
func Read(r io.ReadCloser) (tResponse *TrackerResponse, e error) {
	tr := TrackerResponse{}
	err := bencode.Unmarshal(r, &tr)
	if err != nil {
		// TODO Better error handling
		return &tr, err
	}

	return &tr, err
}

func newRequest(announce string, hash [20]byte, id [20]byte) (query string) {
	base, err := url.Parse(announce)
	if err != nil {
		//log.Errorf("Announce url invalid: %v", url)
	}
	params := url.Values{
		"info_hash": []string{string(hash[:])},
		"peer_id":   []string{string(id[:])},
		"event":     []string{"started"},
		"compact":   []string{"1"},
	}
	base.RawQuery = params.Encode()
	return base.String()

}

func newPeerID() [20]byte {
	var bytes [20]byte
	for i := 0; i < 20; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return bytes
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

/*
Implement a udp connection later
func dial(m) (c *Conn, err error) {
	return net.Dial("udp",
*/
