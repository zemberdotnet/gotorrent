package httpDownload

import (
	"github.com/zemberdotnet/gotorrent/bitfield"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"io/ioutil"
	"net/http"
)

type MirrorConn struct {
	Url     string
	FileExt string
}

func (m MirrorConn) Dial() {
	// Don't keep these connections open so no need to dial now
	return
}

// might be better to return a reader/writer and save on memory?
// sensible max downlods will be needed
func (m MirrorConn) AttemptDownloadPiece(task interfaces.Task) ([]byte, error) {
	req, err := m.buildRequest(task)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func (m MirrorConn) Active() bool {
	return false
}

func (m MirrorConn) Status() int {
	return 0
}

func (m MirrorConn) Bitfield() bitfield.Bitfield {
	return bitfield.Bitfield{}
}

// fill method in
func (m MirrorConn) SetActive() {
	return
}

func (m MirrorConn) SetBitfield(b bitfield.Bitfield) {
	return
}
