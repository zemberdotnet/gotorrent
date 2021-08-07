package connection

import (
	"errors"
	"sync"

	"github.com/zemberdotnet/gotorrent/peer"
)

type ConnectionPool struct {
	Peers       []peer.Peer
	connFactory func(peer.Peer) Connection
	connections chan Connection
	activeConns int
	lock        *sync.Mutex
}

func NewConnectionPool(peers []peer.Peer, factory func(peer.Peer) Connection, maxConn int) *ConnectionPool {
	return &ConnectionPool{
		Peers:       peers,
		connections: make(chan Connection, maxConn),
		activeConns: 0,
		lock:        &sync.Mutex{},
	}
}

// ReturnConnection puts a connections back into the connection pool
func (cp *ConnectionPool) ReturnConnection(conn *PeerConn) {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	cp.connections <- conn
}

// Full returns true if a connection pool is full
func (cp *ConnectionPool) Full() bool {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	return len(cp.connections) == cap(cp.connections)
}

// TODO Rethink/Remove
// AddConnection attemps to add a connection and return true if succesful
func (cp *ConnectionPool) AddConnection(conn *PeerConn) bool {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	if cp.activeConns < cap(cp.connections) {
		return false
	}
	cp.connections <- conn
	return true
}

// Need more testing
// Should next connection block
// NextConnection returns the next available function
func (cp *ConnectionPool) NextConnection() (Connection, error) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	if len(cp.connections) > 0 { // if connections on queue
		return <-cp.connections, nil
	} else if cp.activeConns < cap(cp.connections) { // if more connections can be made
		if len(cp.Peers) < 1 {
			return nil, errors.New("No remaining peers to create connection")
		}
		lastPeer := cp.Peers[len(cp.Peers)-1]
		cp.Peers = cp.Peers[:len(cp.Peers)-1]
		cp.activeConns++
		return cp.connFactory(lastPeer), nil

	} else {
		return nil, errors.New("No available connections")
	}
}

func (cp *ConnectionPool) RemoveConnection(conn *PeerConn) {
	cp.activeConns--
}
