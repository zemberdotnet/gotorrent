package httpDownload

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"io/ioutil"
	"net/http"
	"time"
)

type MirrorConn struct {
	Url     string
	FileExt string
}

// No need to dial first, but once again this is here to
// satisfy the interface. Will likely do some work on the interface
// and this part of the project
func (m MirrorConn) Dial() {
	// Don't keep these connections open so no need to dial now
	return
}

// AttemptDownloadPiece will attempt to download pieces through http-range request
func (m MirrorConn) AttemptDownloadPiece(piece interfaces.Piece) ([]byte, error) {
	req, err := m.buildRequest(piece)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: time.Second * 45,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// TODO: returning the io.Reader may work better
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

// MirrorConn won't use the methods below, but I leave them here to satisfy
// interfaces.Connection. In the future I may change this though as the
// BitTorrent strategy becomes more concrete
func (m MirrorConn) Active() bool {
	return false
}

func (m MirrorConn) Status() int {
	return 0
}

func (m MirrorConn) Bitfield() bitfield.Bitfield {
	return bitfield.Bitfield{}
}

func (m MirrorConn) SetActive() {
	return
}

func (m MirrorConn) SetBitfield(b bitfield.Bitfield) {
	return
}
