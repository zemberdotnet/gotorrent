package client

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	filepath := "/home/snow/Projects/Coding/go/gotorrent/debian.torrent"
	c, err := New(filepath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
}
