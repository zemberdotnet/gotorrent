package interfaces

type Strategy interface {
	Download()
	// These are all communication features
	SetReturnChannel(chan Work)
	RecieveWorkChannel() chan Work
	ReturnWorkChannel() chan Work
	SetPieceChannel(chan interface{})
	PieceChannel() chan interface{}
	// These are properties on the strategy that help create work for it
	Multipiece() bool
	URL() bool
}
