package piece

import (
	"bufio"
	"context"
	"fmt"
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type File struct {
	Pieces    []interfaces.Piece
	pieceChan chan interface{}
	Bitfield  *bitfield.Bitfield
}

// WriteToFile handles the recieving and processing of pieces and then calls
// the private writeToFile when the entire file is downloaded
func (f *File) WriteToFile(end context.CancelFunc, path string) {

	finishedPieces := 0
	// recieving pieces
	go func() {
		for {
			piece := <-f.pieceChan
			// Simple switch...just a framework for the future , but this is good to
			// get started with reflection which may help decouple further in
			// the future
			switch piece.(type) {
			case interfaces.Piece:
				if piece.(interfaces.Piece).Validate() == true {
					if f.containsPiece(piece.(interfaces.Piece)) == false {

						f.Pieces = append(f.Pieces, piece.(interfaces.Piece))
						finishedPieces += piece.(interfaces.Piece).Length()

					} else {
						// this will go later
						fmt.Println("Piece already downloaded")
						f.Bitfield.SetPieceRange(piece.(interfaces.Piece).Index(), piece.(interfaces.Piece).Length())
						continue
					}
					fmt.Printf("Finished %v of %v pieces.\n", finishedPieces, f.Bitfield.Pieces)
				} else {
					fmt.Printf("Piece at index %v failed validation.\n", piece.(interfaces.Piece).Index())
					f.Bitfield.UnsetPieceRange(piece.(interfaces.Piece).Index(), piece.(interfaces.Piece).Length())
				}
			}
		}
	}()

	// waiting for file to finish
	for {
		time.Sleep(time.Second)
		if finishedPieces == f.Bitfield.Pieces {
			end()
			f.writeFileToPath(path)
			break
		}
	}
}

// A simple implementation for now, but in the future we can partially write
// the file as pieces come in so we don't hold the whole file in memory
// Also should balance these with io wear and tear concerns
func (f *File) writeFileToPath(path string) {
	// sorting is quick and makes writing out easier
	sort.Slice(f.Pieces, func(i, j int) bool {
		return f.Pieces[i].Index() > f.Pieces[j].Index()
	})
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed when creating file. Valid path?")
		return
	}

	for _, piece := range f.Pieces {
		reader := bufio.NewReader(piece)
		pc, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println("Failed writing piece to file")
		}
		file.Write(pc)
	}

}

func (f *File) containsPiece(p interfaces.Piece) bool {
	for _, piece := range f.Pieces {
		if p.Index() == piece.Index() {
			if p.Length() == piece.Length() {
				return true
			}
		}
	}
	return false
}

// NewFile creates a new File object
func NewFile(pieceChan chan interface{}, bitfield *bitfield.Bitfield) *File {

	pieces := make([]interfaces.Piece, 0)
	return &File{
		Pieces:    pieces,
		pieceChan: pieceChan,
		Bitfield:  bitfield,
	}
}
