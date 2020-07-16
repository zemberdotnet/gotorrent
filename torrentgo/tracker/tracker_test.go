package tracker

import (
	"fmt"
	//"github.com/zemberdotnet/gotorrent/peer"
	"github.com/zemberdotnet/gotorrent/announce"
	"github.com/zemberdotnet/gotorrent/torrent"
	"testing"
)

func TestGetPeers(t *testing.T) {
	c := make(chan string)
	go announce.Server(c)
	msg := <-c
	fmt.Println("Message:", msg)
	m := &torrent.MetaInfo{
		Announce:  "http://localhost:4000/announce",
		Comment:   "Test Comment",
		Creation:  20000101, // Not sure if this is the proper format
		CreatedBy: "Zember",
		//InfoHash:, // Find proper InfoHash
	}

	tr, err := GetPeers(m)
	if err != nil {
		t.Errorf("Error getting peers: %v", err)
	}
	fmt.Println(tr)
	// fmt.Println(peer.Parse(tr.Peers))
	fmt.Println(msg)
}

func TestHashToString(t *testing.T) {
	var hash = [20]byte{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
	expected := "%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a%0a"
	if hashToString(hash) != expected {
		t.Errorf("Hash differed from expected\n Hash: %v\n Expected:%v\n", hashToString(hash), expected)
	}

}
