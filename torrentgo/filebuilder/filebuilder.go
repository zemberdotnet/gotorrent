package filebuilder

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/piece"
)

type FileBuilder struct {
	pieceCount int
	Pieces     []piece.Piece
	pieceChan  chan *piece.TorrPiece
}

func (f *FileBuilder) Build(end context.CancelFunc, filePath string) error {
	bf := bitfield.NewBitfield(f.pieceCount)
	collectedPieces := 0
	for {
		piece := <-f.pieceChan
		f.Pieces = append(f.Pieces, piece)
		bf.SetPiece(piece.Index())
		collectedPieces++
		fmt.Printf("Percent Downloaded %.2f\n", (float64(collectedPieces)/float64(f.pieceCount))*100)
		if f.pieceCount == collectedPieces {
			break
		}
	}

	f.writeFileToPath(filePath)
	end()
	return nil
}

// WriteToFile handles the recieving and processing of pieces and then calls
// the private writeToFile when the entire file is downloaded
/*
func (f *FileBuilder) WriteToFile(end context.CancelFunc, path string) {
	log.Println("TODO")
	// recieving pieces
	go func() {
		for {
			piece := <-f.pieceChan
			fmt.Println(piece)
			// Simple switch...just a framework for the future , but this is good to
			// get started with reflection which may help decouple further in
			// the future
		}
	}()

	// waiting for file to finish
}
*/

func (f *FileBuilder) Len() int {
	return len(f.Pieces)
}

func (f *FileBuilder) Less(i, j int) bool {
	return f.Pieces[i].Index() < f.Pieces[j].Index()
}

func (f *FileBuilder) Swap(i, j int) {
	f.Pieces[i], f.Pieces[j] = f.Pieces[j], f.Pieces[i]
}

// A simple implementation for now, but in the future we can partially write
// the file as pieces come in so we don't hold the whole file in memory
// Also should balance these with io wear and tear concerns
func (f *FileBuilder) writeFileToPath(path string) {
	// sorting is quick and makes writing out easier
	sort.Slice(f.Pieces, func(i, j int) bool {
		return f.Pieces[i].Index() < f.Pieces[j].Index()
	})

	file, err := os.Create(path)
	if err != nil {
		log.Println("Failed when creating file. Valid path?")
		return
	}

	for _, piece := range f.Pieces {
		reader := bufio.NewReader(piece)
		pc, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Println("Failed writing piece to file")
		}
		file.Write(pc)
	}

}

func (f *FileBuilder) containsPiece(p piece.Piece) bool {
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
func NewFileBuilder(pieceCount int, pieceChan chan *piece.TorrPiece) *FileBuilder {

	pieces := make([]piece.Piece, 0)
	return &FileBuilder{
		Pieces:     pieces,
		pieceChan:  pieceChan,
		pieceCount: pieceCount,
	}
}
