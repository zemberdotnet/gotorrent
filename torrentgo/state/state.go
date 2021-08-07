package state

import (
	"sync"

	"github.com/zemberdotnet/gotorrent/bitfield"
)

type TorrentState struct {
	Length    int
	InfoHash  [20]byte
	Counts    []int
	finished  int
	inProcess int
	lock      *sync.RWMutex
}

func (ts *TorrentState) IncrmentInProcess() {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.inProcess++
}

func (ts *TorrentState) DecrementInProcess() {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	ts.inProcess--
}

func (ts *TorrentState) InProcess() int {
	ts.lock.RLock()
	defer ts.lock.RUnlock()
	return ts.inProcess
}

// does this really need to exits
// data race is inconsequential here
// still would be good practice
func (ts *TorrentState) GetCount(index int) int {
	return ts.Counts[index]
}

func (ts *TorrentState) IncrementCounts(bf bitfield.Bitfield) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	n := len(ts.Counts)

	for i := 0; i < n; i++ {
		if bf.HasPiece(i) {
			ts.Counts[i]++
		}
	}
}

func (ts *TorrentState) DecrementCounts(bf bitfield.Bitfield) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	n := len(ts.Counts)

	for i := 0; i < n; i++ {
		if bf.HasPiece(i) {
			ts.Counts[i]--
		}
	}
}
