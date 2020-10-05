package coordinator

import (
	//	"fmt"
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
)

// A better implementation may be pulling apart the scheduling and
// production of work

type workScheduler interface {
	//Schedule(interfaces.Strategy)
	GetWork() interfaces.Work
	Return(interfaces.Work)
}

type BasicScheduler struct {
	started        bool
	workChan       chan interfaces.Work
	strategyToWork map[interfaces.Strategy]interfaces.Work
	taskFactory    func(bool) interfaces.Task
	connCreator    interfaces.ConnectionCreator
}

func NewBasicScheduler(bt *bitfield.Bitfield, cc interfaces.ConnectionCreator) *BasicScheduler {
	c := make(chan interfaces.Work)
	s := make(map[interfaces.Strategy]interfaces.Work)

	return &BasicScheduler{
		workChan:       c,
		strategyToWork: s,
		taskFactory:    TaskFactory(bt),
		connCreator:    cc,
	}
}

/*
func (b *BasicScheduler) Schedule(s interfaces.Strategy) {
	for _, e := range b.returnedQueue {
		if e.GetStrategy() == s {
			b.workChan <- b.Handle(e)
		}
	}
	if b.started {
		b.workChan <- b.generateWork(s)
	} else {
		b.workChan <- b.generateFirstWork(s)
		b.started = true
	}
}
*/

/*
func (b *BasicScheduler) WorkFactory() {
	// Inspect the
	if len(b.workChan) < 100 {
		b.workChan <- b.generateWork(s)
	}
}
*/
func (b *BasicScheduler) AddStrategyToWork(s interfaces.Strategy, w interfaces.Work) {
	b.strategyToWork[s] = w
}

func (b *BasicScheduler) Handle(w interfaces.Work) interfaces.Work {

	return nil
}

func (b *BasicScheduler) Return(w interfaces.Work) {
	if w.GetTask().Completed() {

	} else {
		// requeu this
	}
}

func (b *BasicScheduler) GetWork() interfaces.Work {
	full := true
	for full {
		for k, _ := range b.strategyToWork {
			//			fmt.Println("LENGTH:", len(k.RecieveWorkChannel()))
			if len(k.RecieveWorkChannel()) < cap(k.RecieveWorkChannel()) {
				full = false
			}
		}
	}

	return b.generateWork()
}

func (b *BasicScheduler) generateWork() interfaces.Work {
	var minStrat interfaces.Strategy
	var minQueueSize int = 100000000
	var _ interfaces.Work
	for k, v := range b.strategyToWork {
		if minStrat == nil {
			minStrat = k
		}

		if len(k.RecieveWorkChannel()) < minQueueSize {
			minQueueSize = len(k.RecieveWorkChannel())
			minStrat = k
			_ = v
		}
	}

	//	return workType.Create()
	return b.NewWorkFromStrategy(minStrat)

}
