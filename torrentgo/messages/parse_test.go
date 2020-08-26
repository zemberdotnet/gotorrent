package messages

import (
	"reflect"
	"testing"
)

func TestParseBasicMessage(t *testing.T) {
	data := []struct {
		input    []byte
		expected *message
		err      error
	}{
		{[]byte{0, 0, 0, 0}, &message{MsgKeepAlive, 0, 0, 0, nil}, nil}, // KeepAliveTest
		// this test case should be updated to work
		{[]byte{1, 1, 1, 1}, &message{MsgKeepAlive, 0, 0, 0, nil}, nil},                       // BadKeepAlive - Should result in an error
		{[]byte{0, 0, 0, 1, MsgChoke}, &message{MsgChoke, 0, 0, 0, nil}, nil},                 // ChokeTest
		{[]byte{0, 0, 0, 1, MsgUnchoke}, &message{MsgUnchoke, 0, 0, 0, nil}, nil},             // MsgUnchoke
		{[]byte{0, 0, 0, 1, MsgInterested}, &message{MsgInterested, 0, 0, 0, nil}, nil},       // InterestedTest
		{[]byte{0, 0, 0, 1, MsgNotInterested}, &message{MsgNotInterested, 0, 0, 0, nil}, nil}, // NotInterestedTest
		{[]byte{0, 0, 0, 1, MsgNotInterested, 1}, nil, ErrInvalidMessageLength},               // Bad input test
	}

	for _, e := range data {
		message, err := parseBasicMessage(e.input)
		if err != e.err {
			t.Errorf("Unexpected or missing error\nActual:\t%v\nExpected:\t%v\n", err, e.err)
		}
		if err != nil {
			t.Logf("Recieved error\nErr:\t%v\n", err)
			continue // return, message will be nil so not safe to continue
		}
		if !reflect.DeepEqual(message, e.expected) {
			t.Errorf("Messages do not match\nActual:\t%v\nExpected:\t%v\n", message, e.expected)
		}
	}
}

func TestParseHaveMessage(t *testing.T) {
	data := []struct {
		input    []byte
		expected *message
		err      error
	}{
		{[]byte{0, 0, 0, 5, MsgHave, 0, 0, 125, 0}, &message{MsgHave, 32000, 0, 0, nil}, nil},
		{[]byte{0, 0, 0, 5, MsgHave, 0, 0, 0, 0}, &message{MsgHave, 0, 0, 0, nil}, nil},
		{[]byte{0, 0, 0, 6, MsgHave, 0, 0, 0, 0}, &message{MsgHave, 0, 0, 0, nil}, ErrInvalidLength},     // Bad length recieved
		{[]byte{0, 0, 0, 5, MsgHave, 255, 255, 255, 255}, &message{MsgHave, 4294967295, 0, 0, nil}, nil}, // Max Uint32 case
		{[]byte{0, 0, 0, 5, MsgHave, 1, 255, 255, 255, 2}, &message{MsgHave, 4294967295, 0, 0, nil}, ErrInvalidMessageLength},
	}
	for _, e := range data {
		message, err := parseHaveMessage(e.input)
		if err != e.err {
			t.Errorf("Unexpected or missing error\nActual:%v\nExpected:%v\n", err, e.err)
		}
		if err != nil {
			t.Logf("Recieved error\nErr:%v\n", err)
			continue // return, message will be nil so not safe to continue
		}
		if !reflect.DeepEqual(message, e.expected) {
			t.Errorf("Messages do not match\nActual:%v\nExpected:%v\n", message, e.expected)
		}
	}
}

func TestParseBitfieldMessage(t *testing.T) {
	data := []struct {
		input    []byte
		expected *message
		err      error
	}{
		{[]byte{0, 0, 0, 5, MsgBitfield, 0, 0, 0, 0}, &message{MsgBitfield, 0, 0, 0, []byte{0, 0, 0, 0}}, nil},
		{[]byte{0, 0, 0, 5, MsgBitfield, 255, 255, 255, 255}, &message{MsgBitfield, 0, 0, 0, []byte{255, 255, 255, 255}}, nil},
		{[]byte{0, 0, 0, 5, MsgBitfield, 255, 1, 0, 255}, &message{MsgBitfield, 0, 0, 0, []byte{255, 1, 0, 255}}, nil},
		{[]byte{0, 0, 0, 2, MsgBitfield, 255}, &message{MsgBitfield, 0, 0, 0, []byte{255}}, nil}, // barely enough to pass
		{[]byte{0, 0, 0, 1, MsgBitfield}, nil, ErrInvalidMessageLength},
		{[]byte{0, 0, 0, 2, MsgBitfield}, nil, ErrInvalidMessageLength},
	}
	for _, e := range data {
		message, err := parseBitfieldMessage(e.input)
		if err != e.err {
			t.Errorf("Unexpected or missing error\nActual:%v\nExpected:%v\n", err, e.err)
		}
		if err != nil {
			t.Logf("Recieved error\nErr:%v\n", err)
			continue // return, message will be nil so not safe to continue
		}
		if !reflect.DeepEqual(message, e.expected) {
			t.Errorf("Messages do not match\nActual:%v\nExpected:%v\n", message, e.expected)
		}

	}

}

func TestParseRequestMessage(t *testing.T) {
	data := []struct {
		input    []byte
		expected *message
		err      error
	}{
		{[]byte{0, 0, 0, 13, MsgRequest, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255}, &message{MsgRequest, 0, 0, 255, nil}, nil},
	}
	for _, e := range data {
		message, err := parseRequestMessage(e.input)
		if err != e.err {
			t.Errorf("Unexpected or missing error\nActual:%v\nExpected:%v\n", err, e.err)
		}
		if err != nil {
			t.Logf("Recieved error\nErr:%v\n", err)
			continue // return, message will be nil so not safe to continue
		}
		if !reflect.DeepEqual(message, e.expected) {
			t.Errorf("Messages do not match\nActual:%v\nExpected:%v\n", message, e.expected)
		}

	}

}

func TestParseCancelMessage(t *testing.T) {
	// Somewhat handled by Request Testing
}

func TestParsePieceMessage(t *testing.T) {
	data := []struct {
		input    []byte
		expected *message
		err      error
	}{
		{[]byte{0, 0, 0, 10, MsgPiece, 0, 0, 0, 1, 0, 0, 0, 1, 233}, &message{MsgPiece, 1, 1, 0, []byte{233}}, nil},
		{[]byte{0, 0, 0, 11, MsgPiece, 0, 0, 0, 1, 0, 0, 0, 1, 233, 0}, &message{MsgPiece, 1, 1, 0, []byte{233, 0}}, nil},
		{[]byte{0, 0, 0, 11, MsgPiece, 0, 0, 0, 1, 0, 0, 0, 1, 233}, nil, ErrInvalidLength},
		{[]byte{0, 0, 0, 9, MsgPiece, 0, 0, 0, 1, 0, 0, 0, 1, 233}, nil, ErrInvalidLength},
		{[]byte{0, 0, 0, 10, MsgPiece, 0, 0, 0, 1, 0, 0, 0, 1}, nil, ErrInvalidMessageLength},
	}
	for _, e := range data {
		message, err := parsePieceMessage(e.input)
		if err != e.err {
			t.Errorf("Unexpected or missing error\nActual:%v\nExpected:%v\n", err, e.err)
		}
		if err != nil {
			t.Logf("Recieved error\nErr:%v\n", err)
			continue // return, message will be nil so not safe to continue
		}
		if !reflect.DeepEqual(message, e.expected) {
			t.Errorf("Messages do not match\nActual:%v\nExpected:%v\n", message, e.expected)
		}
	}
}
