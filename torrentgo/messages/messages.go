package messages

import "fmt"

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

func Read(io.Reader) {
}
