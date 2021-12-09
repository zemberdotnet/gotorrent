package message

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
)

// TODO Are go functions stored as pointers?
// important for messageParser type

// confusing/close names
var ErrInvalidMessageLength error = errors.New("Invalid Message Length")

// could be broken into other erros if needed but this will generally work
// captures the most general case
var ErrInvalidPayload error = errors.New("Invalid Message Payload") // Bad n

var ErrInvalidLength error = errors.New("Invalid Length Recieved")
var ErrInvalidMessageID error = errors.New("Invalid Message ID Recieved")

/*
Wouldn't this be fun
type Error string

func (e Error) Error() string {
	return string(e)
}
*/
var MessageParser = NewMessageParser()

//var funcs map[int]func([]byte) (*message,error) = NewMessageParser()

// use this type?
type messageParser map[int]func(int, []byte) (*Message, error)

/*
type messageParser struct {
	funcs map[int]func([]byte) (*message, error)
}

// is this conncurently safe?
// its just reads so I think so
func NewMessageParser() *messageParser {
	return &messageParser{
		funcs: genMessageHandlers(),
	}
}
*/
func NewMessageParser() messageParser {
	return genMessageHandlers()
}

func (m *Message) ReadFrom(r io.Reader) (n int64, err error) {
	return 0, nil
}

/*
func (mp messageParser) ReadMessage(conn net.Conn) (*Message, error) {

	msg, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	return MessageParser.ParseMessage(msg)

}
*/

// Parse Message takes in bytes, reads the messageID and calls the parsing function
func (mp messageParser) ParseMessage(msgLen int, b []byte) (m *Message, err error) {
	var parser func(int, []byte) (*Message, error)
	parser, ok := mp[int(b[0])]
	if !ok {
		return nil, ErrInvalidMessageID
	}
	// if len is only one then pass in an empty slice
	if len(b) == 1 {
		m, err = parser(msgLen, b[0:0])
	} else {
		m, err = parser(msgLen, b[1:])
	}

	if err != nil {
		return nil, err
	}
	return m, nil
}

// We are entering dangerous teritory with this, but I think it makes things
// Very elegant. This will need _extensive_ testing
// This function could look like this
//func GenerateMessageHandlers() map[int]func([]byte) *message
// Or it could look like this
//func GenreateMessageHandlers()  map[int]func(*message) thing

// Error handling
// Not adding handling for invalid IDs because those should be filtered at highest level
/*
func parseBasicMessage(msgLen int, b []byte) (m *Message, err error) {
	if len(b) == 0 {
		m = &Message{
			MessageID: MsgKeepAlive,
		}
	} else if len(b) == 1 {
		m = &Message{
			MessageID: int(b[0]),
		}
	} else {
		return nil, ErrInvalidMessageLength
	}
	return
}
*/

func parseChokeMessage(msgLen int, b []byte) (m *Message, err error) {
	m = &Message{
		MessageID: MsgChoke,
	}
	return
}

func parseUnchokeMessaage(msgLen int, b []byte) (m *Message, err error) {
	m = &Message{
		MessageID: MsgUnchoke,
	}
	return
}

func parseInterestedMessage(msgLen int, b []byte) (m *Message, err error) {
	m = &Message{
		MessageID: MsgInterested,
	}
	return
}

func parseNotInterestedMessage(msgLen int, b []byte) (m *Message, err error) {
	m = &Message{
		MessageID: MsgNotInterested,
	}
	return
}

func parseHaveMessage(msgLen int, b []byte) (m *Message, err error) {
	if len(b) != 4 {
		log.Println("Error in Have Message")
		return nil, ErrInvalidMessageLength
	}

	m = &Message{
		MessageID: MsgHave,
		Index:     int(binary.BigEndian.Uint32(b[0:4])),
	}
	return
}

// at a higher level we need to check if bitfield matches expected length
// TODO add length parameter to all messages
func parseBitfieldMessage(msgLen int, b []byte) (m *Message, err error) {
	// 4 byte length, message id, 1 piece
	// These check were not right, we parse length seperately
	// TODO: Add checks for length bitfield = length torrent file
	// no need to check on id we do that above
	log.Println("Bitfield message length: ", msgLen)
	log.Println("Bitfiel Size:", len(b))
	m = &Message{
		MessageID: MsgBitfield,
		Payload:   b,
	}
	return

}

func parsePieceMessage(msgLen int, payload []byte) (m *Message, err error) {
	// length, message id, index, begin, block
	if len(payload)+1 != msgLen {
		log.Println("Invalid payload message length")
		return nil, ErrInvalidLength
	}

	m = &Message{
		MessageID: MsgPiece,
		Index:     int(binary.BigEndian.Uint32(payload[0:4])),
		Begin:     int(binary.BigEndian.Uint32(payload[4:8])),
		Payload:   payload[8:],
	}
	return

}

// Don't know if named variables are best here since we are actually init
// in helper function
func parseRequestMessage(msgLen int, b []byte) (*Message, error) {
	if len(b) != 12 {
		log.Println("Request Message Error")
		return nil, ErrInvalidMessageLength
	}

	return &Message{
		MessageID: MsgRequest,
		Index:     int(binary.BigEndian.Uint32(b[0:4])),
		Begin:     int(binary.BigEndian.Uint32(b[4:8])),
		Length:    int(binary.BigEndian.Uint32(b[8:12])),
	}, nil
}

func parseCancelMessage(msgLen int, b []byte) (*Message, error) {
	if len(b) != 12 {
		log.Println("Cancel Message Error")
		return nil, ErrInvalidMessageLength
	}

	return &Message{
		MessageID: MsgCancel,
		Index:     int(binary.BigEndian.Uint32(b[0:4])),
		Begin:     int(binary.BigEndian.Uint32(b[4:8])),
		Length:    int(binary.BigEndian.Uint32(b[8:12])),
	}, nil
}

func genMessageHandlers() map[int]func(int, []byte) (*Message, error) {
	return map[int]func(int, []byte) (*Message, error){
		MsgChoke:         parseChokeMessage,
		MsgUnchoke:       parseUnchokeMessaage,
		MsgInterested:    parseInterestedMessage,
		MsgNotInterested: parseNotInterestedMessage,
		MsgHave:          parseHaveMessage,
		MsgBitfield:      parseBitfieldMessage,
		MsgRequest:       parseRequestMessage,
		MsgPiece:         parsePieceMessage,
		MsgCancel:        parseCancelMessage,
	}
}
