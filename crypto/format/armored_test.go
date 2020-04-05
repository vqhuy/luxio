package format

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"testing"
)

func TestArmor(t *testing.T) {
	buf := &bytes.Buffer{}
	w := ArmoredWriter(buf)
	plain := make([]byte, 64)
	if _, err := rand.Read(plain); err != nil {
		t.Error(err)
	}
	if _, err := w.Write(plain); err != nil {
		t.Error(err)
	}
	if err := w.Close(); err != nil {
		t.Error(err)
	}
	r := ArmoredReader(buf)
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(out, plain) {
		t.Error("decoded value doesn't match")
	}
}
