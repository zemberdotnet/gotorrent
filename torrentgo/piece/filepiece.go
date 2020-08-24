package main

import ()

type Piece interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Validate([]byte) bool // assuming we pass the piece directly
	Offset() int
}

// A file piece is made up of blocks

type FilePiece struct {
	piece []byte
	index int
}

// Collection of all FilePieces
type File struct {
	Pieces []*FilePiece
}

// Write to FilePiece
func (f *FilePiece) Write(b []byte) (n int, err error) {

}

// Reads from FilePiece
func (f *FilePiece) Read(b []byte) (n int, err error) {

}

func (f *FilePiece) Offset() int {
	return f.offset
}

func (f *FilePiece) Validate(b []byte) bool {
	// Code for if it matches the hash the validate it
}

// Concurrently writing file
func (f *File) AssembleAndWrite(out string) error {
}
