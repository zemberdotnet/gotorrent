package peer

import (
	"fmt"
	"net"
	"testing"
)

func TestString(t *testing.T) {
	testIP := net.ParseIP("1.0.0.127")
	var testUint uint16 = 1000

	p := Peer{
		testIP,
		testUint,
	}

	fmt.Println(testIP)
	fmt.Println(testUint)
	fmt.Println(p.String())
}
