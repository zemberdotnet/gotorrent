package piece

import ()

type Piece interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Validate([]byte) bool // assuming we pass the piece directly
	Offset() int
}

type FilePiece struct {
	piece []byte
	index int
	begin int
}

// Write to FilePiece
func (f *FilePiece) Write(b []byte) (n int, err error) {
	return 0, nil

}

// Reads from FilePiece
func (f *FilePiece) Read(b []byte) (n int, err error) {
	return 0, nil

}

func (f *FilePiece) Offset() int {
	return f.index
}

func (f *FilePiece) Validate(b []byte) bool {
	// Code for if it matches the hash the validate it
	return false
}

// Concurrently writing file
func (f *File) AssembleAndWrite(out string) error {
	return nil
}
