package message

import (
	"io"
)

// Needs configurations
// Making these fields but leave them empty could casuse issues in the future
// Thinking of alternatives
type Message struct {
	MessageID int
	Index     int
	Begin     int
	Length    int
	Payload   []byte // bitfield or piece

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

func ReadMessage(r io.Reader) (msg *Message, err error) {
	return nil, nil

}
