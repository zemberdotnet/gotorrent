package piece

import (
	"fmt"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"reflect"
)

type File struct {
	Pieces    []interfaces.Piece
	pieceChan chan interface{}
}

func (f *File) WriteToFile(path string) {
	go func() {
		for {
			piece := <-f.pieceChan
			switch piece.(type) {
			case interfaces.Piece:
				f.Pieces = append(f.Pieces, piece.(interfaces.Piece))
			default:
				fmt.Printf("TYPE: %v\n", reflect.TypeOf(piece))
			}

		}
	}()
}

func NewFile(length int, pieceChan chan interface{}) *File {
	pieces := make([]interfaces.Piece, 0)
	return &File{
		Pieces:    pieces,
		pieceChan: pieceChan,
	}
}
