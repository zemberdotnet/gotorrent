package coordinator

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
)

var Index int = 0

type TaskCreator struct {
	taskChannel chan interfaces.Task
}

func NewTaskCreator() TaskCreator {
	return TaskCreator{}
}

func (t TaskCreator) CreateInitalTasks(bf bitfield.Bitfield, n int) {

}

// Create the next task and see if we can accept multiple pieces
func (t TaskCreator) CreateNextTask(multipiece bool) {

}
