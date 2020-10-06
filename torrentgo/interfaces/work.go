package interfaces

import ()

type Work interface {
	SetStrategy(Strategy)
	SetPiece(Piece) // ad
	GetPiece() Piece
	GetStrategy() Strategy
	GetConnection() Connection
	Do() // do the work
	Completed() bool
}

/*
type BitTorrentWork struct {
	strategy strategy
	task     task
	Conn     Connection
	status   bool
}

type MirrorWork struct {
	strategy strategy
	task     task
	Conn     Connection
	status   bool
}

func NewMirrorWorkGenerator(s strategy) func() *MirrorWork {
	return func() *MirrorWork {
		return &MirrorWork{
			strategy: s,
		}
	}
}

func (m *MirrorWork) setStrategy(s strategy) {
	m.strategy = s
}

func (m *MirrorWork) setConnection(c Connection) {
	m.Conn = c
}

func (m *MirrorWork) setTask(t task) {
	m.task = t
}

func (m *MirrorWork) Do() {
	m.strategy.Download()
	m.strategy.RecieveWorkChannel() <- m
}

func (m *MirrorWork) getStrategy() strategy {
	return m.strategy
}
*/
