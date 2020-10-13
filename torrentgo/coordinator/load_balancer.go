package coordinator

import (
	"github.com/zemberdotnet/gotorrent/interfaces"
	"sync"
)

var lock = sync.RWMutex{}

// A loadBalancer helps keep active strategies to the ideal amount
// While custom maxes can be set, BitTorrent strategies perform best around 30 connections
// WebSeeding requests will perform better a low numbers due to higher overhead per connection
type loadBalancer interface {
	AddStrategy(interfaces.Strategy, int)
	NextStrategy() interfaces.Strategy
	isFree(interfaces.Strategy) bool // will likely remove this in the next version as it is no longer needed
	increaseCount(interfaces.Strategy)
	decreaseCount(interfaces.Strategy)
}

// basicLoadBalancer implements the loadBalancer interface
type basicLoadBalancer struct {
	loads map[interfaces.Strategy]*load
}

// load represents an active strategies current count and max
type load struct {
	count int
	max   int
}

func NewBasicLoadBalancer() *basicLoadBalancer {
	loads := make(map[interfaces.Strategy]*load)
	return &basicLoadBalancer{
		loads: loads,
	}
}

// AddStrategy adds a strategy to the loadbalancer to monitor
func (b *basicLoadBalancer) AddStrategy(s interfaces.Strategy, max int) {
	b.loads[s] = &load{
		count: 0,
		max:   max,
	}
}

// NextStrategy gets the next available strategy
func (b *basicLoadBalancer) NextStrategy() interfaces.Strategy {
	for strat, _ := range b.loads {
		if b.isFree(strat) {
			return strat
		}
	}
	return nil
}

// isFree tells if a strategy is free to be initiated
func (b *basicLoadBalancer) isFree(s interfaces.Strategy) bool {
	load, ok := b.loads[s]
	if !ok {
		return false
	}
	return load.count < load.max
}

// increaseCount adds +1 to the count of a strategy
func (b *basicLoadBalancer) increaseCount(s interfaces.Strategy) {
	lock.Lock()
	defer lock.Unlock()
	load, ok := b.loads[s]
	if !ok {
		return
	}
	load.count++
}

// decreaseCount adds +1 to the count of a strategy
func (b *basicLoadBalancer) decreaseCount(s interfaces.Strategy) {
	lock.Lock()
	defer lock.Unlock()
	load, ok := b.loads[s]
	if !ok {
		return
	}
	load.count--
}
