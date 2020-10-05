package coordinator

import (
	"fmt"
	"testing"
)

func TestInterface(t *testing.T) {
	r := NewBasicLoadBalancer()
	printLoadBalancer(r)
}

func printLoadBalancer(lb loadBalancer) {
	fmt.Println(lb)
}
