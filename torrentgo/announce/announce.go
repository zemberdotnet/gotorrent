package announce

import (
	"fmt"
	//"github.com/zemberdotnet/peer"
	"net/http"
)

// Respond with TrackerResponse to a proper request
type TrackerResponse struct {
	Failure  string
	Interval int
	//Peers    []peer.Peer
}

type ChanHolder struct {
	Server *http.Server
	c      chan string
}

func (ch *ChanHolder) torrentReq(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("event") != "started" {
		fmt.Println("Event not started")
	}
	if r.URL.Query().Get("compact") == "0" {
		fmt.Println("Uncompact response requested")
	}

	// Learn here about go channels
	//TODO: Write this method

	ch.c <- "Hello World!"

}

// This might need to return an error for testing or we may use a channel
func Server(c chan string) {
	srv := &http.Server{}
	ch := &ChanHolder{
		Server: srv,
		c:      c,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/announce", ch.torrentReq)
	srv = &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
