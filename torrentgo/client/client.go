package client

import (
	//"bytes"
	//"fmt"
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/coordinator"
	"github.com/zemberdotnet/gotorrent/piece"
	//"github.com/zemberdotnet/gotorrent/handshake"
	"github.com/zemberdotnet/gotorrent/httpDownload"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/p2p"
	"github.com/zemberdotnet/gotorrent/peer"
	"github.com/zemberdotnet/gotorrent/torrent"
	"github.com/zemberdotnet/gotorrent/tracker"
	//"io"
	"log"
	//"net"
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
		log.Printf("Error getting peers: %v\n, Tracker Response: %v\nProceeding...", err, tr)
	}

	client.MetaInfo = m
	if tr != nil {
		client.Tracker = tr
		client.Peers = tr.Parsed
	}
	return &client, err
}

/*
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
*/

func (c *Client) Coordinate() {
	pieceChan := make(chan interface{})
	file := piece.NewFile(c.MetaInfo.Length(), pieceChan)

	var cc interfaces.ConnectionCreator
	var strats []interfaces.Strategy
	b := bitfield.NewBitfield( /*Number of Pieces*/ c.MetaInfo.Length())
	// this is probably to literal
	if c.Peers != nil && c.MetaInfo.URLList != nil { // torrnet and webseed
		cc = coordinator.NewConnectionFactory(c.MetaInfo.URLList, c.Peers, c.MetaInfo.Info.Name)
		mirrorDL := httpDownload.NewMirrorDownload(c.MetaInfo.Length(), c.MetaInfo.PieceLength())
		torrentDL := p2p.NewTorrentDownload(c.MetaInfo.Length())
		strats = []interfaces.Strategy{mirrorDL, torrentDL}
	} else if c.Peers != nil { // Just torrentDL
		cc = coordinator.NewConnectionFactory(nil, c.Peers, c.MetaInfo.Info.Name)
		torrentDL := p2p.NewTorrentDownload(c.MetaInfo.Length())
		strats = []interfaces.Strategy{torrentDL}
	} else {
		cc = coordinator.NewConnectionFactory(c.MetaInfo.URLList, nil, c.MetaInfo.Info.Name)
		mirrorDL := httpDownload.NewMirrorDownload(c.MetaInfo.Length(), c.MetaInfo.PieceLength())
		strats = []interfaces.Strategy{mirrorDL}
	}
	// if urllist available	&& peer list
	// add both strategies
	// else if urllist
	// mirrorstrategy
	// else
	// peer list strategy
	ws := coordinator.NewBasicScheduler(b, cc)
	lb := coordinator.NewBasicLoadBalancer()

	for _, strategy := range strats {
		strategy.SetPieceChannel(pieceChan)
		ws.AddStrategyToWork(strategy, &coordinator.AbstractWork{})
		lb.AddStrategy(strategy, 30)
	}

	coordinator.Coordinate(lb, ws, nil)

	file.WriteToFile("/home/snow/")
}
