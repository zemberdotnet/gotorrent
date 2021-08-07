package handshake

import (
	"io"
)

// Handshake represents the byte array sent as a Handshake to peers
type Handshake struct {
	InfoHash [20]byte
	PeerID   [20]byte
	Pstr     []byte
}

func NewEmptyHandshake() *Handshake {
	return &Handshake{}
}


// NewHandshake creates a new Handshake from a hash and peerID.
// It currently uses a hardcoded Pstr. This can be changed later as
// the client is updated
func NewHandshake(hash [20]byte, peerID [20]byte) (h *Handshake) {
	return &Handshake{
		InfoHash: hash,
		PeerID:   peerID,
		Pstr:     []byte("BitTorrent protocol"),
	}
}

// WriteTo wraps an io.Writer. So this isn't quite right
// It fuffils the io.WriterTo interface
func (h *Handshake) WriteTo(w io.Writer) (n int64, err error) {
	i, err := w.Write(h.serialize())
	if err != nil {
		return int64(i), err
	}
	return int64(i), err
}

// ReadFrom reads from an io.Reader and returns a Handshake and an error
func (h *Handshake) ReadFrom(r io.Reader) (n int64, err error) {
	buf := make([]byte, 1)
	i, err := io.ReadFull(r, buf)
	if err != nil {
		return int64(i), err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], buf[2+int(buf[0]):22+int(buf[0])])
	copy(peerID[:], buf[23+int(buf[0]):])

	h.Pstr = buf[1 : 1+int(buf[0])]
	h.InfoHash = infoHash
	h.PeerID = peerID

	return int64(49 + i), nil
}

/*
func (h *Handshake) Handshake() (conn net.Conn, err error) {
	// Add a timeout
	// Start using this pattern more in Go
	// We initialize variable in funciton signature (conn)
	// Then we can assign a variable to it.
	// We can't do ":=" with it, but rather do "="
	// This will clean things up a bit
	conn, err = net.DialTimeout("udp", h.Peer.String(), 10*time.Second)
	if err != nil {
		return nil, err
	}
	// Veggiedefender uses a conn.SetDeadline here
	// Look into it
	n, err := conn.Write(h.serialize())
	fmt.Printf("Wrote %v bytes to stream", n)
	if err != nil {
		return nil, err
	}
	// TODO FINISH READ


		if err != nil {
			// TODO Error handling!
			fmt.Println("Error Connecting")
			return nil
			// What are the cases
		}
		// Should be broken down more but for now lets try this
		n, err := conn.Write(h.serialize())
		fmt.Println("BYTES WRITTER:", n)
		if err != nil {
			// TODO as always handle error !
			return nil
		}
		return conn

	return
}
*/

// serialize forms the BitTorrent specified byte array for a handshake
func (h *Handshake) serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr)) // hopefully int coerces easily to byte
	index := 1 + copy(buf[1:len(h.Pstr)+1], h.Pstr)
	index += copy(buf[index:index+8], []byte{0, 0, 0, 0, 0, 0, 0, 0})
	index += copy(buf[index:index+20], h.InfoHash[:])
	index += copy(buf[index:index+20], h.PeerID[:])
	return buf
}
