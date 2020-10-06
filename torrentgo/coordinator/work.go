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

/*
func NewMirrorWorkGenerator(s interfaces.Strategy) func() interfaces.Work {
	return func() interfaces.Work {
		return &MirrorWork{
			strategy: s,
		}
	}
}
*/

func (bs *BasicScheduler) NewWorkFromStrategy(s interfaces.Strategy) interfaces.Work {
	if s.Multipiece() {
		return &AbstractWork{
			strategy: s,
			piece:    bs.pieceFactory(true),
			conn:     bs.connCreator,
		}

	} else {
		return &AbstractWork{
			strategy: s,
			piece:    bs.pieceFactory(false),
			conn:     bs.connCreator,
		}
	}
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
	//go a.strategy.Download()
	a.strategy.RecieveWorkChannel() <- a
}

/*
func (m *MirrorWork) SetStrategy(s interfaces.Strategy) {
	m.strategy = s
}

func (m *MirrorWork) SetConnection(c interfaces.Connection) {
	m.Conn = c
}

func (m *MirrorWork) SetTask(t interfaces.Task) {
	m.task = t
}

func (m *MirrorWork) Do() {
	fmt.Println("ABOUTTO CALL")
	go m.strategy.Download()
	fmt.Println("SENDING WORK TO CHANNEL")
	m.strategy.RecieveWorkChannel() <- m
}

func (m *MirrorWork) GetStrategy() interfaces.Strategy {
	return m.strategy
}
*/
