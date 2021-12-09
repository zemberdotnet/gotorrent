package main

import (
	"github.com/zemberdotnet/gotorrent/client"
)

func main() {
	c, err := client.New("../debian.torrent")
	if err != nil {
		panic(err)
	}
	c.Coordinate()
}
