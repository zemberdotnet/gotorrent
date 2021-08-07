package coordinator

import (
	"context"
	"sync"
	"time"

	"github.com/zemberdotnet/gotorrent/interfaces"
)

type Coordinator struct {
	strats struct {
		creators []func() interfaces.Strategy
		next     int
		total    int
		max      int
	}
}

func NewCoordinator(strats []func() interfaces.Strategy) {
	// TODO
}

func (c *Coordinator) Coordinate(ctx context.Context) {
	wg := sync.WaitGroup{}
	// create strategies
	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				// we trade being more dynamic and abstract
				// for simplicity and less effort
				if c.strats.total < c.strats.max {
					// create strategies
					if c.strats.next < len(c.strats.creators) {
						strat := c.strats.creators[c.strats.next]()
						go strat.Start()
						c.strats.next++
					} else {
						c.strats.next = 0
						strat := c.strats.creators[c.strats.next]()
						go strat.Start()
						c.strats.next++
					}
				}

			}
		}
	}()

	// start and mange pieces
	go func() {
	}()

}

// start piece sending service

// send them download
