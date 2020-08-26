package messages

import (
	"encoding/binary"
	"io"
)

// TODO
func (m *message) WriteTo(w io.Writer) (n int64, err error) {
	return 0, nil
}

// we could map these messages for quicker access
// buildBasicMessage handles keep-alive, choke, unchoke, interested
// and not interested messages. t is the message type
func BuildBasicMessage(t int) []byte {
	if t == MsgKeepAlive {
		return []byte{0, 0, 0, 0}
	} else if t == MsgChoke || t == MsgUnchoke || t == MsgInterested || t == MsgNotInterested {
		return []byte{0, 0, 0, 1, byte(t)}
	}
	// if a incorrect number is passed then passing back nil
	// this isn't safe though so perhaps this isn't the best error handling strategy
	return nil
}

// need work on have suppression
func BuildHaveMessage(index int) []byte {
	message := make([]byte, 9)
	message[3] = 5
	message[4] = MsgHave
	binary.BigEndian.PutUint32(message[5:9], uint32(index))
	// switch to append instead of copy
	return message
}

func BuildBitfieldMessage(b []byte) []byte {
	message := make([]byte, 5+len(b))
	binary.BigEndian.PutUint32(message[0:4], uint32(1+len(b)))
	message[4] = MsgBitfield
	copy(message[5:], b)
	return message
}

func BuildRequestMessage(index, begin, length int) []byte {
	return cancelAndRequestBuilder(index, begin, length)
}

func BuildPieceMessage(index, begin int, payload []byte) []byte {
	message := make([]byte, 13+len(payload))
	binary.BigEndian.PutUint32(message[0:4], uint32(9+len(payload)))
	message[4] = MsgPiece
	binary.BigEndian.PutUint32(message[5:9], uint32(index))
	binary.BigEndian.PutUint32(message[9:13], uint32(begin))
	copy(message[13:], payload)
	return message
}

func BuildCancelMessage(index, begin, length int) []byte {
	return cancelAndRequestBuilder(index, begin, length)
}

func cancelAndRequestBuilder(index, begin, length int) []byte {
	message := make([]byte, 17)
	message[3] = 13
	message[4] = MsgRequest
	binary.BigEndian.PutUint32(message[5:9], uint32(index))
	binary.BigEndian.PutUint32(message[9:13], uint32(begin))
	binary.BigEndian.PutUint32(message[13:17], uint32(length))
	return message
}
