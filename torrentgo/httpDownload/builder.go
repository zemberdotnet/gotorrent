package httpDownload

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

func (m *mirrorPiece) buildRequest(mirror string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", mirror, new(bytes.Buffer))
	if err != nil {
		return
	}
	req.Header.Set("Range", m.buildHeader())
	return
}

func (m *mirrorPiece) buildHeader() string {
	var b1 strings.Builder

	b1.Write([]byte{'b', 'y', 't', 'e', 's', '='})

	b1.WriteString(strconv.Itoa(m.pieceLength * m.offset))

	b1.Write([]byte{'-'})

	ans := m.pieceLength*(m.offset+1) - 1
	if ans > m.fileLength {
		ans = m.fileLength
	}
	b1.WriteString(strconv.Itoa(ans))
	return b1.String()
}
