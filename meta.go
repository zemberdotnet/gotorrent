package main

type TorrentInfo struct {
	Announce  string   `bencode:"announce"`
	Comment   string   `bencode:"comment"`
	Creation  int      `bencode:"creation date"`
	CreatedBy string   `bencode:"created by"`
	Encoding  string   `bencode:"encoding"`
	Info      InfoDict `bencode:"info"`
}

type InfoDict struct {
	Length   int    `bencode:"length"`
	Name     string `bencode:"name"`
	PieceLen int    `bencode:"piece length"`
	Pieces   string `bencode:"pieces"`
}
