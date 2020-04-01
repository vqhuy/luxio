package niceware

import (
	"bytes"
	"strings"
	"testing"
)

func TestSplitToBlocks(t *testing.T) {
	in := []byte{
		0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7,
		0x7, 0x6, 0x5, 0x4, 0x3, 0x2, 0x1, 0x0,
		0x0, 0x1, 0x2, 0x3, 0x0, 0x1, 0x2, 0x3,
		0x7, 0x6, 0x5, 0x4, 0x7, 0x6, 0x5, 0x4,
	}
	want := []byte{0, 1, 1, 0, 0, 1, 1, 1}
	got := splitToBlocks(in)
	if bytes.Equal(got, want) {
		t.Error("convertion")
	}
}

func TestBytesToPassphrase(t *testing.T) {
	bytes := make([]byte, 32)
	words, err := BytesToPassphrase(bytes)
	if err != nil {
		t.Error(err)
	}
	if strings.Join(words[:], "") != "aaaa" {
		t.Error("expect", "aaaa", "got", strings.Join(words[:], ""))
	}
}
