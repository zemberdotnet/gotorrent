package peer

import (
	"strings"
)

type mirror struct {
	mirror  string
	fileExt string
}

func (m *mirror) SetFileExt(fileExt string) {
	m.fileExt = fileExt
}

func CreateMirrorList(mL []string, fileExt string) []*mirror {
	mirrorList := make([]*mirror, len(mL))
	for i, e := range mL {
		mirrorList[i] = &mirror{
			mirror:  e,
			fileExt: fileExt,
		}
	}
	return mirrorList
}

func (m *mirror) String() string {
	if strings.HasSuffix(m.mirror, "/") {
		return m.mirror + m.fileExt
	}
	return m.mirror
}
