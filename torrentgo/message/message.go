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
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return MessageParser.ParseMessage(buf)
}

func SendMessage(w io.Writer, msg *Message) (err error) {
	buf := msg.Serialize()
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
