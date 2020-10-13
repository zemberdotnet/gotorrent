package coordinator

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/httpDownload"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/peer"
	"net"
)

// AbstractConn is a connection that can serve strategies of all types
type AbstractConn struct {
	conn   net.Conn
	bt     bitfield.Bitfield //Adding in Bitfield to manage BitTorrent rarest-first
	status int
	url    bool
	drop   bool
}

// ConnCreator implements the interfaces.ConnectionCreator
type ConnCreator struct {
	URLSource  []string
	PeerSource []peer.Peer
	FileExt    string
	Active     []interfaces.Connection
}

// NewConnectionFactory creates a new connection factory for the strategies to use
func NewConnectionFactory(urlSource []string, peerSource []peer.Peer, FileExt string) interfaces.ConnectionCreator {
	return &ConnCreator{
		URLSource:  urlSource,
		FileExt:    FileExt, // single file only
		PeerSource: peerSource,
	}
}

func (cc *ConnCreator) GetConnection(s interfaces.Strategy) interfaces.Connection {
	if len(cc.Active) == 0 {

		// TODO: Recycle URL connection
		// TODO: If we make this part more robust cases than we can decouple to a greater degree

		if s.URL() {
			if len(cc.URLSource) > 0 {
				var pop string
				cc.URLSource, pop = cc.URLSource[0:len(cc.URLSource)-1], cc.URLSource[len(cc.URLSource)-1]
				return httpDownload.MirrorConn{
					Url:     pop,
					FileExt: cc.FileExt,
				}
			}

		} else {
			// TODO: Create connections for BitTorrent
		}

	} else {
		// TODO: Implement recylced connections
		if s.URL() {

		} else {

		}
	}
	return nil
}

func (cc *ConnCreator) ReturnConnection(interfaces.Connection) {
	//TODO: Return connections to be recylced
}
