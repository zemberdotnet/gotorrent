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

//TODO maybe want to adjust what we return and assign
// i.e. we return 1, error, then we should have one byte in our handshake?

// ReadFrom reads from an io.Reader and returns a Handshake and an error
func (h *Handshake) ReadFrom(r io.Reader) (n int64, err error) {
	// our first byte is the legnth of the pstr
	lenBuf := make([]byte, 1)
	i, err := io.ReadFull(r, lenBuf)
	if err != nil {
		return int64(i), err
	}

	pstrLen := int(lenBuf[0])

	// buf[0] is the length of the pstr and then 48 more bytes for the rest
	handshakeBuf := make([]byte, pstrLen+48)
	i, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return int64(i), err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrLen+8:pstrLen+28])
	copy(peerID[:], handshakeBuf[pstrLen+28:])

	h.Pstr = handshakeBuf[0:pstrLen]
	h.InfoHash = infoHash
	h.PeerID = peerID

	// extra one to account for the length byte
	return int64(i) + 1, nil
}

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
