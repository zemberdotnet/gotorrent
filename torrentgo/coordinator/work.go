package coordinator

import (
	"github.com/zemberdotnet/gotorrent/interfaces"
)

type AbstractWork struct {
	strategy interfaces.Strategy
	piece    interfaces.Piece
	conn     interfaces.ConnectionCreator
	status   bool
}

func (a *AbstractWork) GetPiece() interfaces.Piece {
	return a.piece
}

func (a *AbstractWork) SetStrategy(s interfaces.Strategy) {
	a.strategy = s
}

func (a *AbstractWork) SetPiece(t interfaces.Piece) {
	a.piece = t
}

func (a *AbstractWork) GetStrategy() interfaces.Strategy {
	return a.strategy
}

func (a *AbstractWork) Completed() bool {
	return a.status
}

func (a *AbstractWork) GetConnection() interfaces.Connection {
	return a.conn.GetConnection(a.GetStrategy())
}

func (a *AbstractWork) Do() {
	a.strategy.RecieveWorkChannel() <- a
}
