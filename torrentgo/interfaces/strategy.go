package interfaces

type Strategy interface {
	// These are all communication features
	Start()
	// These are properties on the strategy that help create work for it
	Multipiece() bool
	URL() bool
}

// switch to Download(Piece)
// and Download(p []*Piece)
// remove the work channel
// add method to Piece
