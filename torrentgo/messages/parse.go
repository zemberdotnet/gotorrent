package messages

import (
	"encoding/binary"
	"errors"
	//"fmt"
	"io"
)

// confusing/close names
var ErrInvalidMessageLength error = errors.New("Invalid Message Length")

// could be broken into other erros if needed but this will generally work
// captures the most general case
var ErrInvalidPayload error = errors.New("Invalid Message Payload") // Bad n

var ErrInvalidLength error = errors.New("Invalid Length Recieved")

type messageParser struct {
	funcs map[int]func([]byte) *message
}

func NewMessageParser() *messageParser {
	return &messageParser{
		funcs: genMessageHandlers(),
	}
}

func (m *message) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

// Parse Message takes in bytes, reads the messageID and calls the parsing function
func (mp *messageParser) ParseMessage(b []byte) (m *message, err error) {
	return &message{}, nil
}

// We are entering dangerous teritory with this, but I think it makes things
// Very elegant. This will need _extensive_ testing
// This function could look like this
//func GenerateMessageHandlers() map[int]func([]byte) *message
// Or it could look like this
//func GenreateMessageHandlers()  map[int]func(*message) thing

// Error handling
// Not adding handling for invalid IDs because those should be filtered at highest level
func parseBasicMessage(b []byte) (m *message, err error) {
	if len(b) == 4 {
		m = &message{
			messageID: MsgKeepAlive,
		}

	} else if len(b) == 5 {
		m = &message{
			messageID: int(b[4]),
		}

	} else {
		return nil, ErrInvalidMessageLength
	}
	return
}

func parseHaveMessage(b []byte) (m *message, err error) {
	if len(b) != 9 {
		return nil, ErrInvalidMessageLength
	}

	if b[3] != 5 {
		return nil, ErrInvalidLength
	}

	m = &message{
		messageID: MsgHave,
		index:     int(binary.BigEndian.Uint32(b[5:9])),
	}
	return
}

// at a higher level we need to check if bitfield matches expected length
func parseBitfieldMessage(b []byte) (m *message, err error) {
	// 4 byte length, message id, 1 piece
	if len(b) < 6 {
		return nil, ErrInvalidMessageLength
	}
	length := binary.BigEndian.Uint32(b[0:4])
	if len(b)-4 != int(length) {
		return nil, ErrInvalidLength
	}
	m = &message{
		messageID: MsgBitfield,
		piece:     b[5:],
	}
	return

}

func parsePieceMessage(b []byte) (m *message, err error) {
	// length, message id, index, begin, block
	if len(b) < 14 {
		return nil, ErrInvalidMessageLength
	}
	length := binary.BigEndian.Uint32(b[0:4])
	if len(b)-4 != int(length) {
		return nil, ErrInvalidLength
	}
	m = &message{
		messageID: MsgPiece,
		index:     int(binary.BigEndian.Uint32(b[5:9])),
		begin:     int(binary.BigEndian.Uint32(b[9:13])),
		piece:     b[13:],
	}
	return

}

func parseRequestMessage(b []byte) (*message, error) {
	return cancelAndRequestParser(b)
}

func parseCancelMessage(b []byte) (*message, error) {
	return cancelAndRequestParser(b)
}

func cancelAndRequestParser(b []byte) (m *message, err error) {
	if len(b) != 17 {
		return nil, ErrInvalidMessageLength
	}
	if b[3] != 13 {
		return nil, ErrInvalidLength
	}
	m = &message{
		messageID: int(b[4]),
		index:     int(binary.BigEndian.Uint32(b[5:9])),
		begin:     int(binary.BigEndian.Uint32(b[9:13])),
		length:    int(binary.BigEndian.Uint32(b[13:17])),
	}

	return
}

func genMessageHandlers() map[int]func([]byte) *message {
	m := make(map[int]func([]byte) *message)
	return m
}
