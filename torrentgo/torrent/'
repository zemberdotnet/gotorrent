package torrent

import (
	bencode "github.com/jackpal/bencode-go"
	"log"
	"os"
)

type Torrent struct {
	Announce string `bencode:"announce"`
	Comment  string `bencode:"comment"`
	InfoDict info   `bencode:"info"`
}

type info struct {
	Length   int    `bencode:"length"`
	Name     string `bencode:"name"`
	PeiceLen int    `bencode:"piece length"`
	Peices   string `bencode:"pieces"`
}

func Create(f *os.File) (t *Torrent) {
	data := Torrent{}
	err := bencode.Unmarshal(f, &data)
	if err != nil {
		log.Fatal(err)
	}
	return &data
}
