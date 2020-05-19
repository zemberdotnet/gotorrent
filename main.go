package main

import (
	bencode "github.com/jackpal/bencode-go"
	"log"
	//"net/url"
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

	torrent.NewTracker()
}
