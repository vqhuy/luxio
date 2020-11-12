package engine

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"strconv"
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
	// Password is ... password of length 13, created from the set of
	// case-insensitive alpha-numeric (a-z, 0-9), with a fixed postfix `#A`.
	// For example, `abcdefghij123#A`.
	// This is used in the case that niceware generated passwords are too
	// long for use.
	Password
)

const (
	PinLength   = 6
	PwdLength   = 13
	PwdSizeByte = 32
)

const (
	postFix1     = "#1"
	postFix2     = "#A"
	sep          = "-"
	alphanumeric = "abcdefghijklmnopqrstuvwxyz0123456789"
)

var errInvalidSize = errors.New("Invalid size")
var errUnsupported = errors.New("Unsupported format")

func cryptoRandInt(reader io.Reader, max int64) (int64, error) {
	nBig, err := rand.Int(reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}
	return nBig.Int64(), nil
}

func generatePIN(in []byte) (string, error) {
	reader := bytes.NewReader(in)
	pin := ""
	for i := 0; i < PinLength; i++ {
		ind, err := cryptoRandInt(reader, 10)
		if err != nil {
			return "", err
		}
		pin += strconv.FormatInt(ind, 10)
	}
	return strings.Join(strings.Split(pin, ""), sep), nil
}

func generatePassword(in []byte) (string, error) {
	max := int64(len(alphanumeric))
	reader := bytes.NewReader(in)
	token := ""
	for i := 0; i < PwdLength; i++ {
		ind, err := cryptoRandInt(reader, max)
		if err != nil {
			return "", err
		}
		token += string(alphanumeric[ind])
	}
	return token + postFix2, nil
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
		return generatePIN(in)
	case PlainLowerCase:
		return generatePassphrase(in)
	case WithSpecialCharacters:
		word, err := generatePassphrase(in)
		if err != nil {
			return "", err
		}
		return strings.Title(word) + sep + postFix1, nil
	case Password:
		return generatePassword(in)
	}

	return "", nil
}
