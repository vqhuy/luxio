package format

import (
	"encoding/base64"
	"io"
)

func ArmoredWriter(w io.Writer) io.WriteCloser {
	return base64.NewEncoder(base64.StdEncoding.Strict(), w)
}

func ArmoredReader(r io.Reader) io.Reader {
	return base64.NewDecoder(base64.StdEncoding.Strict(), r)
}
