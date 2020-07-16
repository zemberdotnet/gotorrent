package client

import (
	"github.com/zemberdotnet/gotorrent/peer"
	"github.com/zemberdotnet/gotorrent/torrent"
	"github.com/zemberdotnet/gotorrent/tracker"
	"log"
	"os"
)

type Client struct {
	MetaInfo *torrent.MetaInfo
	Peers    *[]peer.Peer
}

func New(filepath string) (c *Client, e error) {
	client := Client{}

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error: %v\n, Filepath: %v", err, filepath)
	}

	m, err := torrent.Unmarshal(f)
	if err != nil {
		log.Fatalf("Error: %v\n, Filepath: %v", err, filepath)
	}

	tr, err := tracker.GetPeers(m)
	if err != nil {
		log.Fatalf("Error getting peers: %v\n, Tracker Response: %v", err, tr)
	}
	// change types that are returned
	client.MetaInfo = m
	client.Peers = &tr.Parsed

	return &client, err
}
