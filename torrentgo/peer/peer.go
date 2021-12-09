package peer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/zemberdotnet/gotorrent/bitfield"
)

// maybe we change this a little
type Peer interface {
	fmt.Stringer
}

type TorrPeer struct {
	IP       net.IP
	Port     uint16
	Bitfield bitfield.Bitfield
}

func ParseTorrPeers(s string) (p []Peer, e error) {
	if len(s)%6 != 0 {
		return nil, errors.New("invalid compact response from tracker")
	}
	Peers := make([]Peer, 0)
	b := []byte(s)
	for i := 0; i < (len(b) / 6); i++ {
		offset := i * 6
		peer := TorrPeer{
			IP: net.IP(b[offset : offset+4]),
			//IP:   net.IPv4(b[i], b[i+1], b[i+2], b[i+3]),
			Port: binary.BigEndian.Uint16([]byte(b[offset+4 : offset+6])),
			//Port: concatenate(b[i+4], b[i+5]),
		}
		Peers = append(Peers, peer)
	}
	return Peers, nil
}

//TODO remove this or try to see how we can make it work like
// binary.BigEndian.Uint16([]byte(b[offset+4 : offset+6]))
func concatenate(x, y byte) uint16 {
	pow := 10
	i := int(x)
	j := int(y)
	for j >= pow {
		pow *= 10
	}
	return uint16(i*pow + j)
}

func (p TorrPeer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}
