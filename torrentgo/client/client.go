package client

import (
	"context"
	"fmt"
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/coordinator"
	"github.com/zemberdotnet/gotorrent/httpDownload"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/p2p"
	"github.com/zemberdotnet/gotorrent/peer"
	"github.com/zemberdotnet/gotorrent/piece"
	"github.com/zemberdotnet/gotorrent/torrent"
	"github.com/zemberdotnet/gotorrent/tracker"
	"log"
	"os"
)

// Client represents one torrent download
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
func (c *Client) CreateConnCreatorFromStrategies(strats []interfaces.Strategy) interfaces.ConnectionCreator {

	var multi bool = false
	var single bool = false
	var cc interfaces.ConnectionCreator
	for _, strat := range strats {
		if strat.Multipiece() == false {
			single = true
		} else {
			multi = true
		}
	}
	// might be better to just send this logic lower and have nil handling
	if multi && single {
		cc = coordinator.NewConnectionFactory(c.MetaInfo.URLList, c.Peers, c.MetaInfo.Info.Name)
	} else if multi {
		cc = coordinator.NewConnectionFactory(c.MetaInfo.URLList, nil, c.MetaInfo.Info.Name)
	} else {
		cc = coordinator.NewConnectionFactory(nil, c.Peers, c.MetaInfo.Info.Name)
	}
	return cc
}

func (c *Client) Coordinate() {

	hashes := c.MetaInfo.PieceHashes()
	fmt.Println(len(hashes))

	pieceChan := make(chan interface{})

	b := bitfield.NewBitfield( /*Number of Pieces*/ c.MetaInfo.Pieces())
	file := piece.NewFile(pieceChan, b)

	strats := c.CreateStrategies()
	cc := c.CreateConnCreatorFromStrategies(strats)

	workChan := make(chan interfaces.Work)

	ws := coordinator.NewBasicScheduler(b, cc, hashes)
	lb := coordinator.NewBasicLoadBalancer()

	for _, strategy := range strats {
		strategy.SetPieceChannel(pieceChan)
		strategy.SetReturnChannel(workChan)
		ws.AddStrategyToWork(strategy, &coordinator.AbstractWork{})
		lb.AddStrategy(strategy, 30)
	}
	// might be good to make a coordinator object
	coordinator := coordinator.NewCoordinator(lb, ws, workChan)

	ctx, cancel := context.WithCancel(context.Background())

	go coordinator.Coordinate(ctx)

	file.WriteToFile(cancel, "/home/snow/testfile")
}

// CreateStrategies is a helper function to create
// available strategies based on tracker response and torrent file
func (c *Client) CreateStrategies() []interfaces.Strategy {
	var strats []interfaces.Strategy
	if c.Peers != nil && c.MetaInfo.URLList != nil { // torrnet and webseed
		mirrorDL := httpDownload.NewMirrorDownload(c.MetaInfo.Pieces(), c.MetaInfo.PieceLength())
		torrentDL := p2p.NewTorrentDownload(c.MetaInfo.Pieces())
		strats = []interfaces.Strategy{mirrorDL, torrentDL}
	} else if c.Peers != nil { // Just torrentDL
		torrentDL := p2p.NewTorrentDownload(c.MetaInfo.Pieces())
		strats = []interfaces.Strategy{torrentDL}
	} else {
		mirrorDL := httpDownload.NewMirrorDownload(c.MetaInfo.Pieces(), c.MetaInfo.PieceLength())
		strats = []interfaces.Strategy{mirrorDL}
	}
	return strats
}
