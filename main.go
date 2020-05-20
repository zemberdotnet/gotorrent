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
		defer resp.Body.Close()
		//fmt.Println(reflect.TypeOf(resp.Body).String())
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		if err != nil {
			log.Fatal(err)
		}
		tResp := parseResponse(resp.Body)
		fmt.Println(tResp)
	}
}
