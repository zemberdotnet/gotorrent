package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerId   [20]byte
}

// Probably bad to handle error handling here and not elevate it
func peerConn(peer Peer) (c net.Conn, e error) {
	//	peerString := //

	conn, err := net.DialTimeout("tcp", peerString(peer), time.Duration(3)*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error connecting to: %v\nError: %v\n", peerString, err)
	}
	fmt.Println(conn)
	return conn, err
}

// Maybe this func doesnt't belong here.
func peerString(peer Peer) string {
	port := strconv.Itoa(peer.Port)
	return peer.IP + ":" + port
}

// createHandshake creates the bytes to be sent in the handshake.
// Handshake will always be 49 + len Pstr long 68
func (t *TorrentInfo) newHandshake() (h *Handshake) {
	ih := t.Info.hash()
	pd := sha1.Sum([]byte("unique string"))
	fmt.Println("HASH:", ih)
	handshake := Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: ih,
		PeerId:   pd,
	}
	return &handshake
}

// This does come from the blog
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	index := 1
	index += copy(buf[index:], h.Pstr)
	index += copy(buf[index:], make([]byte, 8))
	index += copy(buf[index:], h.InfoHash[:])
	index += copy(buf[index:], h.PeerId[:])
	fmt.Println("Serial Handshake: ", buf)
	fmt.Println("LEN HADNSHEK:", len(buf))
	return buf
}

func ReadHandshake(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	pstrlen := int(lengthBuf[0])
	fmt.Println("LEN:", pstrlen)
	if pstrlen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}
	// Should be 48 moved to the whole for now
	handshakeBuf := make([]byte, 48+pstrlen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}
	var infoHash, peerID [20]byte
	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuf[pstrlen+8+20:])
	h := Handshake{
		Pstr:     string(handshakeBuf[0:pstrlen]),
		InfoHash: infoHash,
		PeerId:   peerID,
	}
	return &h, nil
}
