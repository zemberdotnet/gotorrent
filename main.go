package main

import (
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"io"
	"log"
	"net"
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
		return
	}
	defer resp.Body.Close()

	//fmt.Println(reflect.TypeOf(resp.Body).String())

	//	body, err := ioutil.ReadAll(resp.Body)

	//	fmt.Println(string(body))
	//if err != nil {
	//	log.Fatal(err)
	//}
	tResp, err := parseResponse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tResp)
	handshake := torrent.newHandshake()
	hbytes := handshake.Serialize()
	//	buf := make([]byte, 0, 4096)
	//	tmp := make([]byte, 256)
	var conn net.Conn
	for i := 0; i != 100; i++ {
		conn, err = peerConn(tResp.Peers[i])
		if err != nil {
			continue
		}
		break
	}
	defer conn.Close()
	n, err := conn.Write(hbytes)
	//n, err := fmt.Fprint(conn, hbytes)
	fmt.Println(n)
	if err != nil {
		fmt.Println(err)
	}
	h, err := ReadHandshake(conn)

	if err != nil && err != io.EOF {
		fmt.Println(err)
	}
	fmt.Println("PSTR", h.Pstr)
	fmt.Println("InfoHash", h.InfoHash)
	fmt.Println("PeerID", h.PeerId)
	/*
		for {
			n, err := conn.Read(tmp)
			if err != nil {
				if err != io.EOF {
					fmt.Println("read error:", err)
				}
				break
			}
			buf = append(buf, tmp[:n]...)
		}

		fmt.Println("total size:", len(buf))
		fmt.Println(buf[0])
		fmt.Println(string(buf[1:]))
		//fmt.Println(conn)
	*/
}
