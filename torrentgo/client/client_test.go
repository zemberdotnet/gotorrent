package client

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
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

	c.ConnectToPeers()

}

/*
func DontTestDownload(t *testing.T) {
	filepath := "/home/snow/Projects/Coding/go/gotorrent/arch.torrent"
	c, err := New(filepath)
	if err != nil {
		fmt.Println(err)
	}

	c.DownloadHTTP()
}
*/

func DontTestShaSum(t *testing.T) {
	filepath := "/home/snow/Projects/Coding/go/gotorrent/torrentgo/client/test.iso"
	fmt.Println("hello")
	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	b, _ := ioutil.ReadAll(f)
	fmt.Println(sha1.Sum(b))
	fmt.Println("above")

}
