package connection

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/zemberdotnet/gotorrent/peer"
)

//  TestNewPeerConn tests the NewPeerConn function
func TestNewPeerConn(t *testing.T) {

}

func ConnectionWorker(ip string) net.Conn {
	conn, err := net.DialTimeout("tcp", ip, 6*time.Second)
	if err != nil {
		log.Printf("Error connecting to peer: %s", err)
	}
	defer conn.Close()
	fmt.Println("Connected to peer")
	return conn
}

func TestInitialize(t *testing.T) {

	//TODO
	/*
		bf := bitfield.NewBitfield(0)
		infoHash := [20]byte{124, 245, 84, 40, 50, 86, 23, 253, 222, 145, 15, 229, 91, 121, 171, 114, 190, 147, 121, 36}
		peerID := [20]byte{72, 72, 72, 64, 62, 65, 63, 62, 63, 62, 63, 62, 63, 62, 63, 62, 63, 62, 63, 62}
	*/

	peer1 := peer.TorrPeer{
		IP:   net.ParseIP("21.231.234.200"),
		Port: 58113,
	}
	peer2 := peer.TorrPeer{
		IP:   net.ParseIP("189.185.124.240"),
		Port: 15026,
	}

	go func() {
		ip1 := net.JoinHostPort(peer1.IP.String(), strconv.Itoa(int(peer1.Port)))
		ip2 := net.JoinHostPort(peer2.IP.String(), strconv.Itoa(int(peer2.Port)))
		go ConnectionWorker(ip1)
		go ConnectionWorker(ip2)
	}()
	time.Sleep(time.Second * 10)

	/*
		peerConn := NewPeerConn(peer, infoHash, peerID)
		err := peerConn.Initialize()
		if err != nil {
			t.Errorf("Error initializing peer connection: %s", err)
		}
	*/

}

func TestSetBitfield(t *testing.T) {
	t.Errorf("TBD")
}

func TestHandshake(t *testing.T) {
	t.Errorf("TBD")
}

func TestSendMessage(t *testing.T) {
	t.Errorf("TBD")
}

func TestRecieveMessage(t *testing.T) {
	t.Errorf("TBD")
}
