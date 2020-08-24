package httpDownload

import (
	"strings"
)

// FilePiece will be moved to another package in the next few updates
// Working towards dedicated file construction, tracking, and management
type FilePiece struct {
	mirrorPiece *mirrorPiece
	piece       []byte
	status      bool
}

// mirrorPiece represents one piece to request from a mirror
type mirrorPiece struct {
	fileLength  int
	pieceLength int
	offset      int
}

type mirror string

// pieceLength will be variable as we need to fill in variable Length pieces
func MirrorPieceGenerator(fileLength int) func(int, int) *mirrorPiece {
	return func(offset, pieceLength int) *mirrorPiece {
		return &mirrorPiece{
			fileLength:  fileLength,
			pieceLength: pieceLength,
			offset:      offset,
		}
	}
}

func MirrorGenerator(fileExt string) func(string) mirror {
	return func(m string) mirror {
		if strings.HasSuffix(m, fileExt) {
			return mirror(m)
		} else {
			return mirror(m + fileExt)
		}
	}
}

func (m mirror) String() string {
	return string(m)
}
