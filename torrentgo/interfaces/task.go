package interfaces

type Task interface {
	Completed() bool
	Index() int
	Length() int
}
