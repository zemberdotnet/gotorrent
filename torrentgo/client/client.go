package client

import (
	"fmt"
	"github.com/zemberdotnet/gotorrent/handshake"
	"github.com/zemberdotnet/gotorrent/torrent"
	"github.com/zemberdotnet/gotorrent/tracker"
	"log"
	"os"
)

type Client struct {
	InputPath  string
	OutputPath string
	InfoHash   [20]byte
	PeerID     [20]byte
}

func New(inputPath string, outputPath string) (c *Client, e error) {
	client := Client{
		InputPath:  inputPath,
		OutputPath: outputPath,
	}
	return &client, nil
}

// Issac told me to put this here but honestly I don't know why
func openFile(path string) (f *os.File, e error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error Reading File: %v\nError:%v\n", path, err)
		return nil, err
	}
	return file, err
}

func (c *Client) CreateTracker() (tr *tracker.Request) {
	f, err := os.Open(c.InputPath)
	if err != nil {
		log.Printf("Error Reading File: %v\nError:%v\n", c.InputPath, err)
		return nil
	}
	torr := torrent.Create(f)
	return tracker.NewTracker(torr)
}

func (c *Client) RemoveFunc(tr *tracker.Request) {
	c.InfoHash = tr.InfoHash
	c.PeerID = tr.PeerID
}

func (c *Client) GetPeers(t *tracker.Request) (tr *tracker.Response) {
	resp, err := tracker.Get(t)
	if err != nil {
		log.Fatalf("Error fetching tracker: %v\nError: %v\n", t, err)
		return nil
	}
	defer resp.Body.Close()
	trackerResp := tracker.Read(resp.Body)
	return trackerResp
}

func (c *Client) Shake() {
	h := handshake.NewHandshake(c.InfoHash, c.PeerID)
	b := h.Serialize()
	fmt.Println(b)
}
