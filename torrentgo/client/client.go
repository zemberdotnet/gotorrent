package client

import (
	"bytes"
	"fmt"
	"github.com/zemberdotnet/gotorrent/handshake"
	//	"github.com/zemberdotnet/gotorrent/httpDownload"
	"github.com/zemberdotnet/gotorrent/peer"
	"github.com/zemberdotnet/gotorrent/torrent"
	"github.com/zemberdotnet/gotorrent/tracker"
	"io"
	"log"
	"net"
	"os"
)

// Client is due to be rewritten.

type Client struct {
	MetaInfo *torrent.MetaInfo
	Tracker  *tracker.TrackerResponse
	Peers    []peer.Peer
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

	client.MetaInfo = m
	client.Peers = tr.Parsed
	client.Tracker = tr

	return &client, err
}

func (c *Client) ConnectToPeers() {
	var conn net.Conn
	var err error
	if len(c.Peers) != 0 {
		for i := 0; i < 20; i++ {

			fmt.Println("Now connecting to:", c.Peers[i])
			h := handshake.NewHandshake(c.Peers[i], c.MetaInfo.InfoHash, c.Tracker.PeerID)

			conn, err = h.Handshake()
			if err != nil {
				fmt.Println("Error Connecting:", err)
				continue
			} else {
				break
			}
		}
		defer conn.Close()
		var buf bytes.Buffer

		io.Copy(&buf, conn)
		fmt.Println(buf.Len())
	}
}
