package messages

// Needs configurations
// Making these fields but leave them empty could casuse issues in the future
// Thinking of alternatives
type message struct {
	messageID int
	index     int    // index of piece
	begin     int    // block offset
	length    int    // length of requested data
	piece     []byte // bitfield or piece
}

const (
	MsgKeepAlive = iota
	MsgChoke
	MsgUnchoke
	MsgInterested
	MsgNotInterested
	MsgHave
	MsgBitfield
	MsgRequest
	MsgPiece
	MsgCancel
	MsgPort
)
