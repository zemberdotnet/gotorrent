package client

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/zemberdotnet/gotorrent/connection"
	"github.com/zemberdotnet/gotorrent/coordinator"
	"github.com/zemberdotnet/gotorrent/filebuilder"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/p2p"
	"github.com/zemberdotnet/gotorrent/peer"
	"github.com/zemberdotnet/gotorrent/piece"
	"github.com/zemberdotnet/gotorrent/state"
	"github.com/zemberdotnet/gotorrent/torrent"
	"github.com/zemberdotnet/gotorrent/tracker"
	"github.com/zemberdotnet/gotorrent/work"
)

// Client represents one torrent download
type Client struct {
	MetaInfo *torrent.MetaInfo
	Tracker  *tracker.TrackerResponse
	Peers    []peer.Peer
	state    *state.TorrentState
}

func New(filepath string) (c *Client, e error) {
	client := Client{}

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error: %v\n, Filepath: %v", err, filepath)
	}
	defer f.Close()

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
	log.Println(m.InfoHash)

	return &client, err
}

func (c *Client) Coordinate() {

	c.state = state.NewTorrentState(c.MetaInfo.PieceLength(), c.MetaInfo.Length(), c.MetaInfo.Pieces(), c.MetaInfo.PieceHashes(), c.MetaInfo.InfoHash)

	outChan := make(chan *piece.TorrPiece, 30)

	coord := coordinator.NewCoordinator(c.state, c.CreateStrategyCreators(outChan))
	ctx, cancel := context.WithCancel(context.Background())
	/*
		for _, peer := range c.Peers {
			conn, err := net.DialTimeout("tcp", peer.String(), time.Second*2)
			if err != nil {
				log.Println(err)
				continue
			}
			conn.Close()
			log.Println("Success!")

		}
	*/

	fb := filebuilder.NewFileBuilder(c.MetaInfo.Pieces(), outChan)

	go fb.Build(cancel, "./debian")
	go coord.Coordinate(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Second * 3)
		}
	}
}

type CreateStrategyInput struct {
	singlePieceWorkQueue chan piece.Piece
}

func (c *Client) CreateStrategyCreators(outChan chan *piece.TorrPiece) []func() interfaces.Strategy {
	wc := work.NewWorkCreator(c.state)

	pool := connection.NewConnectionPool(c.Peers, connection.DefaultConnectionFactory(c.MetaInfo.InfoHash, c.Tracker.PeerID), 30)
	strats := make([]func() interfaces.Strategy, 0)

	if c.Peers != nil {
		torrentDL := p2p.NewP2PFactory(c.state, pool, wc, outChan)

		strats = append(strats, torrentDL)
	}

	/*
		if c.MetaInfo.URLList != nil {
			mirrorDL := func() interfaces.Strategy {
				return httpDownload.NewMirrorDownload(c.MetaInfo.Pieces(), c.MetaInfo.PieceLength())
			}
			strats = append(strats, mirrorDL)
		}
	*/

	return strats

}

func (c *Client) CreateConnectionPool() *connection.ConnectionPool {

	DefaultConFactory := connection.DefaultConnectionFactory(c.MetaInfo.InfoHash, c.Tracker.PeerID)

	return connection.NewConnectionPool(c.Peers, DefaultConFactory, 30)
}
