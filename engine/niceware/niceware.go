package niceware

import (
	"errors"
)

var errInvalidSize = errors.New("Invalid size")
var errInvalidByte = errors.New("Invalid byte encountered")

func splitToBlocks(in []byte) []byte {
	var bytes []byte
	for i := 0; i < 4; i++ {
		first := byte(0)
		second := byte(0)
		for j := 0; j < 4; j++ {
			first ^= in[i*8+j]
			second ^= in[i*8+4+j]
		}
		bytes = append(bytes, first, second)
	}
	return bytes
}

func BytesToPassphrase(in []byte) ([]string, error) {
	if len(in) != 32 {
		return nil, errInvalidSize
	}

	bytes := splitToBlocks(in)
	var words []string

	for i := 0; i < 8; i += 2 {
		cur := bytes[i]
		next := bytes[i+1]
		wordIndex := int(cur)*256 + int(next)
		word := wordList[wordIndex]
		if word == "" {
			return nil, errInvalidByte
		} else {
			words = append(words, word)
		}
	}
	return words, nil
}
