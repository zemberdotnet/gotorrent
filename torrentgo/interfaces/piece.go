package interfaces

type Piece interface {
	Index() int
	Length() int
	Begin() int
	Validate() bool

	Write([]byte) (int, error)
	Read([]byte) (int, error)
}
