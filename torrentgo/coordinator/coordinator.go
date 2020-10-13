package coordinator

import (
	"context"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"sync"
	"time"
)

type coordinator struct {
	lb loadBalancer
	ws workScheduler
	ch chan interfaces.Work
}

func NewCoordinator(lb loadBalancer, ws workScheduler, d chan interfaces.Work) coordinator {
	return coordinator{
		lb: lb,
		ws: ws,
		ch: d,
	}
}

func (crd coordinator) Coordinate(ctx context.Context) {
	wg := sync.WaitGroup{}
	// Start strats
	wg.Add(3)
	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				strat := crd.lb.NextStrategy()
				if strat != nil {
					go strat.Download() // puts startegy into ready and waiting
					crd.lb.increaseCount(strat)
				}
			}
		}
	}()

	// Create and dispatch work
	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				crd.ws.GetWork()
			}
		}
	}()

	// Return work
	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			select {
			case <-ctx.Done():
				wg.Done()
				return
			default:
				returnedWork := <-crd.ch
				crd.lb.decreaseCount(returnedWork.GetStrategy())
				crd.ws.Return(returnedWork)

			}
		}
	}()

	wg.Wait()
}
