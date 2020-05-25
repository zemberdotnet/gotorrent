package main

import (
	"fmt"
	"github.com/zemberdotnet/gotorrent/client"
	"log"
	"os"
)

func main() {
	input := string(os.Args[1])
	output := string(os.Args[2])

	c, err := client.New(input, output)
	if err != nil {
		log.Fatal(err)
	}

	t := c.CreateTracker()
	c.RemoveFunc(t)
	peers := c.GetPeers(t)

	c.Shake()
	fmt.Println(peers)

}
