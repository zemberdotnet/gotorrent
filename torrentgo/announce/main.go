package announce

import (
	"fmt"
)

func main() {
	c := make(chan string)
	go Server(c)
	msg := <-c
	fmt.Println(msg)
}
