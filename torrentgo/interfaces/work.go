package interfaces

type Work interface {
	SetStrategy(Strategy)
	SetPiece(Piece) // ad
	GetPiece() Piece
	GetStrategy() Strategy
	GetConnection() Connection
	Do() // do the work
	Completed() bool
}
