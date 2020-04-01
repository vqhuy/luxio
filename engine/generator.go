package engine

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/vqhuy/luxio/engine/niceware"
)

const (
	// PIN is a 6-digit PIN code, separated by `-`.
	// For example, `0-1-2-3-4-5`.
	PIN = iota
	// PlainLowerCase is niceware generated 4-word password,
	// separated by `-`.
	// For example, `hello-world-from-luxio`.
	PlainLowerCase
	// WithSpecialCharacters is niceware generated 4-word password,
	// title cased, separated by `-`, with a fixed postfix `#1`.
	// For example, `Hello-World-From-Luxio-#1`.
	WithSpecialCharacters
)

const (
	PwdSizeByte = 32
)

const (
	postFix = "#1"
	sep     = "-"
)

var errInvalidSize = errors.New("Invalid size")
var errUnsupported = errors.New("Unsupported format")

// generatePIN splits a 32-byte byte array to 4 blocks and XORs
// them together, and converts the result to an `Uint64`, which will
// be used as a seed for `math/rand`.
func generatePIN(in []byte) string {
	bytes := make([]byte, 8)
	for j := 0; j < 8; j++ {
		for i := 0; i < 4; i++ {
			bytes[j] ^= in[i*8+j]
		}
	}
	seed := binary.LittleEndian.Uint64(bytes)
	s := rand.NewSource(int64(seed))
	r := rand.New(s)
	pin := fmt.Sprintf("%06d", r.Intn(1000000))
	return strings.Join(strings.Split(pin, ""), sep)
}

func generatePassphrase(in []byte) (string, error) {
	words, err := niceware.BytesToPassphrase(in)
	if err != nil {
		return "", err
	}
	return strings.Join(words, sep), nil
}

func Generate(in []byte, format int) (string, error) {
	if len(in) != PwdSizeByte {
		return "", errInvalidSize
	}
	switch format {
	case PIN:
		return generatePIN(in), nil
	case PlainLowerCase:
		return generatePassphrase(in)
	case WithSpecialCharacters:
		word, err := generatePassphrase(in)
		if err != nil {
			return "", err
		}
		return strings.Title(word) + sep + postFix, nil
	}

	return "", nil
}
