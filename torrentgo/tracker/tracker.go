package tracker

import (
	"io"

	bencode "github.com/jackpal/bencode-go"
	"github.com/zemberdotnet/gotorrent/peer"
	"github.com/zemberdotnet/gotorrent/torrent"

	"errors"
	"math/rand"
	"net/http"
	"net/url"
)

// TODO: Implement UDP connection
//       Clean eror handling
//       Improving naming

// Tracker Response represents a bencoded response from a tracker
type TrackerResponse struct {
	Failure  string `bencode:"failure"`
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
	Parsed   []peer.Peer
	PeerID   [20]byte
}

// GetPeers
func GetPeers(m *torrent.MetaInfo) (t *TrackerResponse, e error) {
	t = &TrackerResponse{}
	id := newPeerID()
	if m.GetAnnounce() == "" {
		return nil, errors.New("No announce url or hash")
	}
	url := newRequest(m.GetAnnounce(), m.GetHash(), id)
	t.PeerID = id
	return t.getPeers(url)
}

func (t *TrackerResponse) getPeers(url string) (*TrackerResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return &TrackerResponse{}, err
	}
	tResp, err := t.Read(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	return tResp, nil
}

// Format neutral to deal with bad responses from tracker
func (t *TrackerResponse) Read(r io.ReadCloser) (tResponse *TrackerResponse, e error) {
	err := bencode.Unmarshal(r, t)
	if err != nil {
		// TODO Better error handling
		return t, err
	}

	t.Parsed, err = peer.ParseTorrPeers(t.Peers)
	if err != nil {

		return nil, err
		// TODO Handle Error
	}

	return t, err
}

func newRequest(announce string, hash [20]byte, id [20]byte) (query string) {
	base, err := url.Parse(announce)
	if err != nil {
		// Handle Error
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
