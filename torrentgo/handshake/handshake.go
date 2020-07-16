package handshake

import (
	"net"
)

type Handshake struct {
	InfoHash [20]byte
	PeerID   [20]byte
	Pstr     []byte
}

func New(hash [20]byte, peerID [20]byte) (h *Handshake) {
	return &Handshake{
		InfoHash: hash,
		PeerID:   peerID,
		Pstr:     []byte("BitTorrent protocol"),
	}
}

func (h *Handshake) Handshake() (conn *net.Conn) {
	return nil
}

func (h *Handshake) serialize() []byte {
	// use copy here
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr)) // hopefully int coerces easily to byte
	index := 1 + copy(buf[1:len(h.Pstr)+1], h.Pstr)
	index += copy(buf[index:index+8], []byte{0, 0, 0, 0, 0, 0, 0, 0})
	index += copy(buf[index:index+20], h.InfoHash[:])
	index += copy(buf[index:index+20], h.PeerID[:])
	// We haven't handeled errors here
	return buf
}
