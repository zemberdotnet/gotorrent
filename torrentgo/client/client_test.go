package client

import (
	//	"crypto/sha1"
	"fmt"
	//"io/ioutil"
	//"os"
	"testing"
)

func TestNew(t *testing.T) {

	filepath := "/home/snow/Projects/Coding/go/gotorrent/debian.torrent"
	c, err := New(filepath)
	if err != nil {
		fmt.Println(err)
	}

	if len(c.Peers) == 0 {
		t.Errorf("Failed to retrieve peers")
	}

	// fmt.Println(c)
	// fmt.Println(c.Peers)

	//c.ConnectToPeers()

}

func TestClient(t *testing.T) {
	filepath := "/home/snow/Projects/Coding/go/gotorrent/debian.torrent"
	c, err := New(filepath)
	if err != nil {
		t.Errorf("ERROR: %v", err)
	}

	c.Coordinate()
}
