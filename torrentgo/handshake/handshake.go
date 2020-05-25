package handshake

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

// Dial Connection
// Send Overconnection
// Read From Connection

func NewHandshake(i [20]byte, p [20]byte) *Handshake {
	pstr := "BitTorrent protocol"
	h := Handshake{
		Pstr:     pstr,
		InfoHash: i,
		PeerID:   p,
	}
	return &h
}

func (h *Handshake) Serialize() []byte {
	b := make([]byte, 68)
	index := 1
	b[0] = byte(19)
	index += copy(b[index:], []byte(h.Pstr))
	index += copy(b[index:], make([]byte, 8))
	index += copy(b[index:], h.InfoHash[:])
	index += copy(b[index:], h.PeerID[:])
	return b
}

//func Dial(p *tracker.Peer) (conn *net.Conn, e error) {
//}
