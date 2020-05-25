package main

import (
	"fmt"
	"net/url"
)

func main() {
	ih := [5]byte{49, 49, 49, 49, 49}
	params := url.Values{
		"info_hash": []string{string(ih[:])},
	}
	fmt.Println(params.Encode())
}
