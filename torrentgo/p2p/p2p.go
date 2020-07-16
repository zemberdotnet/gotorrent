package p2p

import (
	"github.com/zemberdotnet/gotorrent/peer"
	"net"
)

type Connection struct {
	PeerID   [20]byte
	InfoHash [20]byte
	Peer     peer.Peer
	Choked   bool
	Conn     *net.Conn
}
