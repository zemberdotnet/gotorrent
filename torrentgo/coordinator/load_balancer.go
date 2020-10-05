package coordinator

import (
	"fmt"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"sync"
)

var lock = sync.RWMutex{}

type loadBalancer interface {
	AddStrategy(interfaces.Strategy, int)
	NextStrategy() interfaces.Strategy
	isFree(interfaces.Strategy) bool
	increaseCount(interfaces.Strategy)
	decreaseCount(interfaces.Strategy)
}

type load struct {
	count int
	max   int
}

// pointer helps keep things consisten across rw
// I think...
type basicLoadBalancer struct {
	loads map[interfaces.Strategy]*load
}

func NewBasicLoadBalancer() *basicLoadBalancer {
	loads := make(map[interfaces.Strategy]*load)
	return &basicLoadBalancer{
		loads: loads,
	}
}

func (b *basicLoadBalancer) AddStrategy(s interfaces.Strategy, max int) {
	b.loads[s] = &load{
		count: 0,
		max:   max,
	}
}

func (b *basicLoadBalancer) NextStrategy() interfaces.Strategy {
	// could be placed elsewhere more efficiently
	for k, _ := range b.loads {
		if b.isFree(k) {
			return k
		}
	}
	return nil
}

// Adding a RLock is actually bad here since it prohibits the regular Lock from activating
// isFree tells if a strategy is free to be initiated
func (b *basicLoadBalancer) isFree(s interfaces.Strategy) bool {
	l, ok := b.loads[s]
	if !ok {
		return false
	}
	return l.count < l.max
}

// increaseCount adds +1 to the count of a strategy
func (b *basicLoadBalancer) increaseCount(s interfaces.Strategy) {
	lock.Lock()
	defer lock.Unlock()
	l, ok := b.loads[s]
	if !ok {
		return
	}
	l.count++
	fmt.Println(l.count)
}

// decreaseCount adds +1 to the count of a strategy
func (b *basicLoadBalancer) decreaseCount(s interfaces.Strategy) {
	lock.Lock()
	defer lock.Unlock()
	l, ok := b.loads[s]
	if !ok {
		return
	}
	l.count--
}
