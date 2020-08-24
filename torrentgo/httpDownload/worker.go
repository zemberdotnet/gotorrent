package httpDownload

import (
	"errors"
	"net/http"
	"time"
)

// At a certain point validate should fufill an interface as validating a connection and validating
// a mirror share some common properties
func (m mirror) validate() error {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Head(m.String())
	if err != nil {
		return err
	}
	if len(resp.Header["Accept-Ranges"]) != 0 && resp.Header["Accept-Ranges"][0] != "bytes" {
		return errors.New("Mirror does not accept ranges")
	}
	// While this could fit in, I think it would be better to validate this in real time
	// if len(resp.Header["Content-Length"]) != 0 && resp.Header["Content-Length"][0] != strconv.Itoa(m.mirrorPiece.fileLength) {
	//	return errors.New("Content length does not match")
	//}
	return nil
}
