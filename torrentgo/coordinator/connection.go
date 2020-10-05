package coordinator

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/httpDownload"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/peer"
	"net"
)

type AbstractConn struct {
	conn   net.Conn
	bt     bitfield.Bitfield
	status int
	url    bool
	drop   bool
}

type ConnCreator struct {
	URLSource  []string
	PeerSource []peer.Peer
	Active     []interfaces.Connection
}

func NewConnectionFactory(urlSource []string, peerSource []peer.Peer) interfaces.ConnectionCreator {
	return &ConnCreator{
		URLSource:  urlSource,
		PeerSource: peerSource,
	}
}

func (cc *ConnCreator) GetConnection(s interfaces.Strategy) interfaces.Connection {
	if len(cc.Active) == 0 {
		// should probably recycle these connections
		if s.URL() {
			// it will break down here if this is nil so we should work on this
			// we can improve by not copying / setting up a better queue
			if len(cc.URLSource) > 0 {
				var pop string
				cc.URLSource, pop = cc.URLSource[0:len(cc.URLSource)-1], cc.URLSource[len(cc.URLSource)-1]
				return httpDownload.MirrorConn{
					Url:     pop,
					FileExt: "",
				}
			}
			// Create using URL strategy

		} else {
			// Create using BitTorrent strategy

		}

	} else {
		if s.URL() {

		} else {

		}
	}
	return nil
}

func (cc *ConnCreator) ReturnConnection(interfaces.Connection) {

}
