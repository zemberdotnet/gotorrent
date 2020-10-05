package interfaces

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
)

type ConnectionCreator interface {
	GetConnection(Strategy) Connection
	ReturnConnection(Connection)
}

type Connection interface {
	Dial()
	AttemptDownloadPiece(Task) ([]byte, error)
	Active() bool
	Status() int
	Bitfield() bitfield.Bitfield
	SetActive()
	SetBitfield(bitfield.Bitfield)
}

// type hosts will keep the connection alive
// if it has a active connection then when dial is called it will return
// the active connection
