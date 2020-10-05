package coordinator

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
)

type abstractTask struct {
	index     int
	begin     int
	length    int
	completed bool
}

func (a abstractTask) Type() string {
	return "task"
}

func (a abstractTask) Completed() bool {
	return a.completed
}

func (a abstractTask) Index() int {
	return a.index
}

func (a abstractTask) Length() int {
	return a.length
}

func TaskFactory(bt *bitfield.Bitfield) func(multi bool) interfaces.Task {
	return func(multi bool) interfaces.Task {
		if multi {
			index, length := bt.LargestGap()
			return abstractTask{
				index:  index,
				length: length,
			}

		} else {
			// need to add in handling for single pieces
			return abstractTask{
				index:  0,
				length: 0, // make this the bittorent size
			}

		}

	}
}
