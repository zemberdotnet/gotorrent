package httpDownload

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/zemberdotnet/gotorrent/interfaces"
)

// buildRequest creates a http range request to download pieces
func (m MirrorConn) buildRequest(task interfaces.Piece) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", m.Url+m.FileExt, new(bytes.Buffer))
	if err != nil {
		return
	}
	req.Header.Set("Range", m.buildHeader(task))

	return
}

func (m MirrorConn) buildHeader(task interfaces.Piece) string {
	var b1 strings.Builder

	b1.Write([]byte{'b', 'y', 't', 'e', 's', '='})

	b1.WriteString(strconv.Itoa(PieceLength * task.Index()))

	b1.Write([]byte{'-'})
	ans := (task.Length()*PieceLength + PieceLength*task.Index()) - 1
	b1.WriteString(strconv.Itoa(ans))
	return b1.String()
}
