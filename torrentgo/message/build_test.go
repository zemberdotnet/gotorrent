package message

import (
	"bytes"
	"testing"
)

func TestBuildBasicMessage(t *testing.T) {
	data := []struct {
		message  int
		length   int
		expected []byte
	}{
		{0, 4, []byte{0, 0, 0, 0}},
		{1, 5, []byte{0, 0, 0, 1, 1}},
		{2, 5, []byte{0, 0, 0, 1, 2}},
		{3, 5, []byte{0, 0, 0, 1, 3}},
		{4, 5, []byte{0, 0, 0, 1, 4}},
	}
	for _, e := range data {
		message := BuildBasicMessage(e.message)
		if len(message) != e.length || !bytes.Equal(message, e.expected) {
			t.Errorf("Basic message build failed.\nActual:\t%v\nExpeted:\t%v\n", message, e.expected)
		}
	}
}

func TestBuildHaveMessage(t *testing.T) {
	data := []struct {
		index    int
		expected []byte
	}{
		{0, []byte{0, 0, 0, 5, MsgHave, 0, 0, 0, 0}},
		{32000, []byte{0, 0, 0, 5, MsgHave, 0, 0, 125, 0}},
	}
	for _, e := range data {
		message := BuildHaveMessage(e.index)
		if !bytes.Equal(message, e.expected) {
			t.Errorf("Have message failed.\nActual:\t%v\nExpected:\t%v\n", message, e.expected)
		}
	}
}

func TestBuildBitfieldMessage(t *testing.T) {
	data := []struct {
		bitfield []byte
		expected []byte
	}{
		{[]byte{1, 1}, []byte{0, 0, 0, 3, MsgBitfield, 1, 1}},
		{[]byte{225, 225, 225}, []byte{0, 0, 0, 4, MsgBitfield, 225, 225, 225}},
	}
	for _, e := range data {
		message := BuildBitfieldMessage(e.bitfield)
		if !bytes.Equal(message, e.expected) {
			t.Errorf("Bitfield message failed.\nActual:\t%v\nExpected:\t%v\n", message, e.expected)
		}
	}

}

func TestBuildRequestMessage(t *testing.T) {
	// constraints - length should be 2^14 or 2^15 bytes
	data := []struct {
		index    int
		begin    int
		length   int
		expected []byte
	}{
		{0, 0, 32000, []byte{0, 0, 0, 13, MsgRequest, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 125, 0}},
	}
	for _, e := range data {
		message := BuildRequestMessage(e.index, e.begin, e.length)
		if !bytes.Equal(message, e.expected) {
			t.Errorf("Request message differs from expected\nActual:\t%v\nExpected:\t%v\n", message, e.expected)

		}
	}
}

func TestBuildPieceMessasge(t *testing.T) {
	data := []struct {
		index    int
		begin    int
		payload  []byte
		expected []byte
	}{
		{0, 0, []byte{1}, []byte{0, 0, 0, 10, MsgPiece, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{32000, 32000, []byte{1, 1}, []byte{0, 0, 0, 11, MsgPiece, 0, 0, 125, 0, 0, 0, 125, 0, 1, 1}},
	}
	for _, e := range data {
		message := BuildPieceMessage(e.index, e.begin, e.payload)
		if !bytes.Equal(message, e.expected) {
			t.Errorf("Request message differs from expected\nActual:\t%v\nExpected:\t%v\n", message, e.expected)
		}
	}
}

func TestBuildCancelMessage(t *testing.T) {
	// constraints - length should be 2^14 or 2^15 bytes
	data := []struct {
		index    int
		begin    int
		length   int
		expected []byte
	}{
		{0, 0, 32000, []byte{0, 0, 0, 13, byte(MsgRequest), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 125, 0}},
	}
	for _, e := range data {
		message := BuildCancelMessage(e.index, e.begin, e.length)
		if !bytes.Equal(message, e.expected) {
			t.Errorf("Request message differs from expected\nActual:\t%v\nExpected:\t%v\n", message, e.expected)

		}
	}
}
