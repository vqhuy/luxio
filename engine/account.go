package engine

import (
	"net/url"
	"strings"

	"golang.org/x/crypto/cryptobyte"
)

type AccountInfo struct {
	Domain   string
	Username string
}

// Serialize returns length-prefixed, byte array from account info.
func (acc *AccountInfo) Serialize() []byte {
	b := &cryptobyte.Builder{}
	for _, in := range [][]byte{
		[]byte(acc.Domain), []byte(acc.Username),
	} {
		b.AddUint16LengthPrefixed(func(b *cryptobyte.Builder) {
			b.AddBytes(in)
		})
	}
	return b.BytesOrPanic()
}

func (acc *AccountInfo) Deserialize(in []byte) (string, string) {
	input := cryptobyte.String(in)
	var result []string
	var value cryptobyte.String
	for !input.Empty() {
		if !input.ReadUint16LengthPrefixed(&value) {
			panic("bad format")
		}
		result = append(result, string(value))
	}
	return result[0], result[1]
}

func DomainFromUrl(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	if u.Scheme == "" {
		return link
	}
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) <= 2 {
		return u.Hostname()
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}
