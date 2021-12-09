package interfaces

import "context"

type Strategy interface {
	// These are all communication features
	Start(context.Context)
	// These are properties on the strategy that help create work for it
}

// switch to Download(Piece)
// and Download(p []*Piece)
// remove the work channel
// add method to Piece
