package main

import (
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	filepath := string(os.Args[1])
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error Reading File\n Error:", err)
	}
	torrent := TorrentInfo{}
	err = bencode.Unmarshal(f, &torrent)

	tr, err := torrent.NewTracker()
	resp, err := tr.getReq()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("We made it")
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}
