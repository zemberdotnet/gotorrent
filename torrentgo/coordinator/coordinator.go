package coordinator

import (
	"context"
	"sync"
	"time"

	"github.com/zemberdotnet/gotorrent/interfaces"
	"github.com/zemberdotnet/gotorrent/state"
)

type Coordinator struct {
	strats struct {
		creators []func() interfaces.Strategy
		next     int
		total    int
		max      int
	}
	state *state.TorrentState
}

func NewCoordinator(state *state.TorrentState, creators []func() interfaces.Strategy) *Coordinator {
	// TODO make this configurable
	max := 30

	return &Coordinator{
		strats: struct {
			creators []func() interfaces.Strategy
			next     int
			total    int
			max      int
		}{
			creators: creators,
			next:     0,
			total:    0,
			max:      max,
		},
		state: state,
	}
}

func (c *Coordinator) Coordinate(ctx context.Context) {
	count := 0
	wg := sync.WaitGroup{}
	sleepDuration := 25 * time.Millisecond
	// create strategies
	go func() {
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				// if we have less strats going than allowed then start more
				if c.strats.total < c.strats.max {
					count++
					// if the next pointer on strats is less than the legnth
					// use it to start otherwise take the pointer back to zero
					// and iterate through the strats again
					if c.strats.next < len(c.strats.creators) {
						strat := c.strats.creators[c.strats.next]()
						go strat.Start(ctx)
						c.strats.total++
						c.strats.next++
					} else {
						c.strats.next = 0
						strat := c.strats.creators[c.strats.next]()
						go strat.Start(ctx)
						c.strats.total++
						c.strats.next++
					}
				} else {
					time.Sleep(sleepDuration)
				}
			}
		}
	}()
}

// start piece sending service

// send them download
