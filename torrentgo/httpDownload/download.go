package httpDownload

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// Download is under heavy construction. I do not plan to retain this file as it is, but
// it served as a Proof of Concept. It succesfully downloaded files under good conditions.
// Now I plan to break this file into other packages/parts and make everything more elegant

type download struct {
	maxWorkers  int
	fileLength  int
	pieceLength int
	numOfPieces int // Could change this implementation
}

func NewDownload(maxWorkers, fileLength, numOfPieces int) download {
	return download{
		maxWorkers:  maxWorkers,
		fileLength:  fileLength,
		numOfPieces: numOfPieces,
	}
}

// Considerations:
// Making a buffer of size d.fileLength is expensive. Future versions
// Will switch to buffered io and a dedicated file forming package

// Very crude, but succesfully downloads as a Proof of Concept
func (d download) DownloadFromMirrors(mirrors []string, fileExt string) []byte {
	// Hard code 30 XD
	buf := make([]byte, d.fileLength)
	c := make(chan *FilePiece, d.numOfPieces)
	workers := 0
	i := 0
	queue := d.GenMirrorPieces()
	for {
		if time.Now().Second()%10 == 0 {
			fmt.Println(len(queue), workers)

		}

		if len(queue) > 0 {
			if workers < d.maxWorkers {
				m := queue[0].NewMirror(mirrors[i] + fileExt)
				i++
				go m.download(c)
				queue = queue[1:]
				workers++
			}
		}

		if len(c) > 0 {
			workers--
			temp := <-c
			if temp.status {
				copy(buf[temp.mirrorPiece.offset*temp.mirrorPiece.pieceLength:], temp.piece)
			} else {
				queue = append(queue, temp.mirrorPiece)
			}
		}
		if len(queue) == 0 && workers == 0 {
			fmt.Println("breaKing")
			break
		}

	}
	return buf
}

// Also very crude. MirrorReq data type has been restructed and this function
// will be tossed/changed
func (m *mirrorReq) download(c chan<- *FilePiece) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(m.req)
	if err != nil {
		c <- &FilePiece{
			mirrorPiece: m.mirrorPiece,
		}
		return
	} else if resp.StatusCode != 206 {
		c <- &FilePiece{
			mirrorPiece: m.mirrorPiece,
		}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c <- &FilePiece{
			mirrorPiece: m.mirrorPiece,
		}
		return
	}
	c <- &FilePiece{
		status:      true,
		mirrorPiece: m.mirrorPiece,
		piece:       body,
	}
}

// Generator for mirror pieces using the mirror piece generator function and callback
func (d download) GenMirrorPieces() []*mirrorPiece {
	// num of pieces, dynamic later, hard now
	mPieces := make([]*mirrorPiece, d.numOfPieces)
	// bad idea... yes it is :-)
	pieceLength := int(math.Ceil(float64(d.fileLength) / float64(d.numOfPieces)))
	generator := MirrorPieceGenerator(d.fileLength, pieceLength)
	for i := 0; i < d.numOfPieces; i++ {
		mPieces[i] = generator(i)
	}
	return mPieces
}
