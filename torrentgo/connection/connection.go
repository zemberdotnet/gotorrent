package connection

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/handshake"
	"github.com/zemberdotnet/gotorrent/message"
	"github.com/zemberdotnet/gotorrent/peer"
)

type Connection interface {
}

// PeerConn represents a connection to a Peer
type PeerConn struct {
	Conn       net.Conn
	Choked     bool
	Interested bool
	Bitfield   bitfield.Bitfield
	Peer       peer.Peer
	infoHash   [20]byte
	peerID     [20]byte
	Backlog    int
	Err        error
}

func DefaultConnectionFactory(infoHash, peerID [20]byte) func(peer.Peer) *PeerConn {
	return func(peer peer.Peer) *PeerConn {
		return NewPeerConn(peer, infoHash, peerID)
	}
}

// NewPeerConn creates a new __unconnected__ PeerConn
func NewPeerConn(peer peer.Peer, infoHash, peerID [20]byte) *PeerConn {
	return &PeerConn{
		Conn:     nil,
		Choked:   true,
		Peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}
}

func (p *PeerConn) Initialize() error {
	log.Println("Initializing connection to peer:", p.Peer.String())
	conn, err := net.DialTimeout("tcp", p.Peer.String(), 20*time.Second)
	if err != nil {
		log.Println("Error connecting to peer:", err)
		p.Err = err
		return err
	}
	log.Println("Connected to peer:", p.Peer.String())
	p.Conn = conn

	// we do not need the response at this time so have
	_, err = p.Handshake()
	if err != nil {
		log.Println("Error handshaking with peer:", err)
		p.Err = err
		p.Conn.Close()
		return err
	}

	err = p.SetBitfield()
	if err != nil {
		log.Println("Error setting bitfield:", err)
		p.Err = err
		p.Conn.Close()
		return err
	}

	// TODO MORE MORE MORE

	return nil
}

// TODO TODO TODO
func (p *PeerConn) SetBitfield() error {
	p.Conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer p.Conn.SetDeadline(time.Time{})

	msg, err := message.ReadMessage(p.Conn)
	if err != nil {
		log.Println("Error reading message from peer:", err)
		return err
	}
	log.Printf("Got Bitfield Message: Message Type: %s\n", msg.String())

	if msg.MessageID != message.MsgBitfield {
		err := fmt.Errorf("expected bitfield message but got %d", msg.MessageID)
		log.Println(err)
		return err
	}

	p.Bitfield = bitfield.NewBitfieldFromBytes(msg.Payload)

	return err
}

func (p *PeerConn) Handshake() (*handshake.Handshake, error) {
	// set Deadline to three seconds then at the end
	//  we set connetion to not timeout so we can keep alive
	p.Conn.SetDeadline(time.Now().Add(time.Second * 3))
	defer p.Conn.SetDeadline(time.Time{})

	// create a new handshake using infohash and peerID
	// then write handshake to the connection
	hs := handshake.NewHandshake(p.infoHash, p.peerID)
	_, err := hs.WriteTo(p.Conn)
	if err != nil {
		log.Println("Error writing handshake to peer:", err)
		return nil, err
	}

	// new handshake for response
	resp := handshake.NewEmptyHandshake()
	_, err = resp.ReadFrom(p.Conn)
	if err != nil {
		log.Println("Error reading handshake from peer:", err)
		// TODO
		return nil, err
	}

	// TODO Implment check on bytes equality in infohash

	return resp, err
}

func (p *PeerConn) SendMessage(m *message.Message) error {
	return nil
}

func (p *PeerConn) ReadMessage() (*message.Message, error) {
	return nil, nil
}
