package coordinator

import (
	"fmt"
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/httpDownload"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"time"
)

/*
type BasicCoordinator struct {
	bitfields []bitfield.Bitfield
}
*/

// Maybe need something that is a method of the bitfield, like
// take in bitfield -> process -> use bitfield to determine next best worker

func Coordinate(lb loadBalancer, ws workScheduler, d chan interfaces.Work) {
	// Worker
	// Workers should maintain a connection until that connection fails
	sent := 0
	failed := 0

	// Start strats
	// No way to stop
	go func() {
		for {
			strat := lb.NextStrategy()
			if strat != nil {
				go strat.Download() // puts startegy into ready and waiting
				lb.increaseCount(strat)
				fmt.Println("started a strat")
			}
		}
	}()

	// Create and dispatch work
	go func() {
		for {
			work := ws.GetWork()
			go work.Do()
			sent++
		}
	}()

	// Return work
	go func() {
		for {
			returnedWork := <-d
			lb.decreaseCount(returnedWork.GetStrategy())
			ws.Return(returnedWork)

		}
	}()
	time.Sleep(time.Second * 6)
	fmt.Println(sent, failed)
}

func mainCoordinate() {
	// create bitfield
	bytearray := make([]byte, 100000)
	bf := bitfield.Bitfield{
		Bitfield: bytearray,
	}
	strat := httpDownload.NewMirrorDownload(123456789, 1234)
	d := strat.ReturnWorkChannel()

	cc := NewConnectionFactory([]string{"http://example.com"}, nil)

	ws := NewBasicScheduler(&bf, cc)
	ws.AddStrategyToWork(strat, &AbstractWork{})
	lb := NewBasicLoadBalancer()
	lb.AddStrategy(strat, 10)

	Coordinate(lb, ws, d)
}
