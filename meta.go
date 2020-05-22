package main

// TorrentInfo defines the structure of a .torrent file in Go objects
// We use bencode-go to unmarshal .torrent files into the TorrentInfo struct
type TorrentInfo struct {
	Announce  string   `bencode:"announce"`
	Comment   string   `bencode:"comment"`
	Creation  int      `bencode:"creation date"`
	CreatedBy string   `bencode:"created by"`
	Encoding  string   `bencode:"encoding"`
	Info      InfoDict `bencode:"info"`
}

// InfoDict represents the info key and its subvalues from the .torrent file.
// It is a substructure of TorrentInfo
type InfoDict struct {
	Length   int    `bencode:"length"`
	Name     string `bencode:"name"`
	PieceLen int    `bencode:"piece length"`
	Pieces   string `bencode:"pieces"`
}
