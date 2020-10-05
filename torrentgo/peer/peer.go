package peer

import (
	"errors"
	"github.com/zemberdotnet/gotorrent/bitfield"
	"net"
	"strconv"
)

type Peer struct {
	IP       net.IP
	Port     uint16
	Bitfield bitfield.Bitfield
}

func ParsePeers(s string) (p []Peer, e error) {
	if len(s)%6 != 0 {
		return nil, errors.New("Invalid compact response from tracker")
	}
	Peers := make([]Peer, 0)
	b := []byte(s)
	for i := 0; i < len(b)-6; i++ {
		peer := Peer{
			IP:   net.IPv4(b[i], b[i+1], b[i+2], b[i+3]),
			Port: concatenate(b[i+4], b[i+5]),
		}
		Peers = append(Peers, peer)
	}
	return Peers, nil
}

func concatenate(x, y byte) uint16 {
	pow := 10
	i := int(x)
	j := int(y)
	for j >= pow {
		pow *= 10
	}
	return uint16(i*pow + j)
}

func (p Peer) String() string {
	return p.IP.String() + ":" + strconv.Itoa(int(p.Port))
}
