package message

import (
	"encoding/binary"
	"io"
	"log"
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
	MsgChoke = iota
	MsgUnchoke
	MsgInterested
	MsgNotInterested
	MsgHave
	MsgBitfield
	MsgRequest
	MsgPiece
	MsgCancel
	MsgPort
	MsgKeepAlive
)

func ReadMessage(r io.Reader) (msg *Message, err error) {
	// Can't use io ReadAll this blocks and blows up!
	msgLenBuf := make([]byte, 4)
	_, err = io.ReadFull(r, msgLenBuf)
	if err != nil {
		return nil, err
	}

	msgLen := binary.BigEndian.Uint32(msgLenBuf)
	// KeepAlive messgea
	if msgLen == 0 {
		return &Message{MessageID: MsgKeepAlive}, nil
	}

	msgBuf := make([]byte, msgLen)
	// what happens on a zero byte read
	_, err = io.ReadFull(r, msgBuf)
	if err != nil {
		return nil, err
	}

	msg, err = MessageParser.ParseMessage(int(msgLen), msgBuf)
	if err != nil {
		return nil, err
	}
	return msg, nil

}

func SendMessage(w io.Writer, msg *Message) (err error) {
	log.Printf("Sending Message type: %s\n", msg.String())
	buf := msg.Serialize()
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (msg *Message) String() string {
	switch msg.MessageID {
	case MsgKeepAlive:
		return "KeepAlive"
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitfield:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgCancel:
		return "Cancel"
	case MsgPort:
		return "Port"
	}
	return "Unknown"
}
