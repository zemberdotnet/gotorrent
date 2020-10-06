package httpDownload

import (
	"bytes"
	"fmt"
	"github.com/zemberdotnet/gotorrent/interfaces"
	"net/http"
	"strconv"
	"strings"
)

// probaly would be moved to another package or not who know
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

	b1.WriteString(strconv.Itoa(task.Length() * task.Index()))

	b1.Write([]byte{'-'})

	ans := task.Length()*(task.Index()+1) - 1
	b1.WriteString(strconv.Itoa(ans))
	fmt.Println(b1.String())
	return b1.String()
}
