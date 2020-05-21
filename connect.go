package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Handshake struct {
	Pstrlen  byte
	Pstr     []byte
	Reserved []byte
}

// Probably bad to handle error handling here and not elevate it
func peerConn(peer Peer) (c net.Conn) {
	//	peerString := //

	conn, err := net.DialTimeout("tcp", peerString(peer), time.Duration(3)*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error connecting to: %v\nError: %v\n", peerString, err)
	}
	fmt.Println(conn)
	return conn
}

// Maybe this func doesnt't belong here.
func peerString(peer Peer) string {
	port := strconv.Itoa(peer.Port)
	return peer.IP + ":" + port
}

/*
func initHandshake (conn net.Conn)
	Pstr = []byte("BitTorrent protocol")
	Pstrlen = byte('19')
	Reserved = []byte("00000000")
	n, p
}
*/
