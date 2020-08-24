package handshake

import (
	"fmt"
	"net"
	//"github.com/google/go-cmp/cmp"
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	//	"testing/iotest"
	"time"
)

const clientID = "BitTorrent protocol"

// serializetestFunction mocks the h.serialize() method and lets us test index
func serializeTestFunction(h *Handshake, expectedIndex int) func(t *testing.T) {
	return func(t *testing.T) {
		buf := make([]byte, len(h.Pstr)+49)
		buf[0] = byte(len(h.Pstr)) // hopefully int coerces easily to byte
		index := 1 + copy(buf[1:len(h.Pstr)+1], h.Pstr)
		index += copy(buf[index:index+8], []byte{0, 0, 0, 0, 0, 0, 0, 0})
		index += copy(buf[index:index+20], h.InfoHash[:])
		index += copy(buf[index:index+20], h.PeerID[:])
		if index != expectedIndex {
			t.Errorf("Ending index differed from Expected")
		}
	}
}

func TestSerialize(t *testing.T) {
	// fuzzing random test cases
	for i := 0; i < 100000; i++ {
		h := &Handshake{
			InfoHash: genRandomArray(),
			PeerID:   genRandomArray(),
			Pstr:     []byte("BitTorrent protocol"),
		}
		s := h.serialize()
		if s[0] != 19 {
			t.Errorf("Length of pstr miscalcuated")
		}
		if !bytes.Equal(s[1:20], []byte(clientID)) {
			t.Errorf("Pstr does not match expected\nPstr: %v\n Expected:%v\n", s[1:20], []byte(clientID))
		}
		if !bytes.Equal(s[20:28], []byte{0, 0, 0, 0, 0, 0, 0, 0}) {
			t.Errorf("Reserved bytes not set correctly actual:%v\n", s[20:28])
		}
		if !bytes.Equal(s[28:48], h.InfoHash[:]) {
			t.Errorf("InfoHash does not match exected\n Actual:\t%v\nExpected:\t%v\n", s[28:48], h.InfoHash)
		}
		if !bytes.Equal(s[48:], h.PeerID[:]) {
			t.Errorf("PeerID does not match expected\n Actual:\t%v\nExpected:\t%v\n", s[48:], h.PeerID)
		}

	}
	// Testing index
	for j := 0; j < 100000; j++ {
		h := &Handshake{
			InfoHash: genRandomArray(),
			PeerID:   genRandomArray(),
			Pstr:     []byte("BitTorrent protocol"),
		}
		t.Run("fuzzed-case-index", serializeTestFunction(h, 68))
	}
	// Testing empty input
	h := &Handshake{}
	if len(h.serialize()) != 49 {
		t.Errorf("Genereated invalid output from empty array")
	}
}

func DontTestWriteTo(t *testing.T) {
	go func() {
		conn, err := net.Dial("tcp", ":8000")
		if err != nil {
			t.Errorf("Failed to start connection")
		}
		defer conn.Close()
		conn.Write([]byte("Hello"))
	}()

	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		t.Errorf("Failed to start connection")
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			t.Errorf("Failed to read from connection")
			return
		}
		fmt.Println(string(buf))
		return
	}

}

func TestWriteTo(t *testing.T) {
	infoHash := genRandomArray()
	peerID := genRandomArray()
	h := &Handshake{
		InfoHash: infoHash,
		PeerID:   peerID,
		Pstr:     []byte("BitTorrent protocol"),
	}

	buf := new(bytes.Buffer)
	h.WriteTo(buf)
	if buf.Bytes()[0] != 19 {
		t.Errorf("Pstrlen differs from expected")
	}
	fmt.Println(string(buf.Bytes()))

}

func TestWriteToConnection(t *testing.T) {
	var _ io.ReaderFrom = (*Handshake)(nil)
	c := make(chan bool)
	go func() {
		l, err := net.Listen("tcp", ":8000")
		if err != nil {
			t.Errorf("Failed to listen to connection")
			return
		}
		c <- true
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				t.Errorf("Failed to accept connection")
				return
			}
			defer conn.Close()
			buf, err := ioutil.ReadAll(conn)
			if err != nil {
				t.Errorf("Failed to read connection")
				return
			}
			fmt.Println(string(buf))
			return
		}
	}()
	// Blocking until we are ready to dial
	// Not sure we need assignment
	_ = <-c

	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		t.Errorf("Failed to read connection\n%v", err)
		return
	}
	defer conn.Close()
	h := &Handshake{
		InfoHash: genRandomArray(),
		PeerID:   genRandomArray(),
		Pstr:     []byte("BitTorrent protocol"),
	}
	h.WriteTo(conn)
}

func genRandomArray() [20]byte {
	rand.Seed(time.Now().UnixNano())
	//rand.Seed(127)
	var buf [20]byte
	for i := 0; i < 20; i++ {
		v := rand.Intn(127)
		buf[i] = byte(v)
	}
	return buf
}
