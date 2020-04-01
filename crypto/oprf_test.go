package crypto

import (
	"bytes"
	"crypto/sha512"
	"testing"

	"github.com/vqhuy/luxio/engine"
	"golang.org/x/crypto/hkdf"
)

func prf(pwd string, ctx *engine.AccountInfo, dev *deviceForAccount) []byte {
	pwdPoint := newRistrettoPoint()
	h := hkdf.New(sha512.New, []byte(pwd), nil, ctx.Serialize())
	b := make([]byte, 64)
	h.Read(b)
	pwdPoint.FromUniformBytes(b)

	sa := newRistrettoKey()
	sa.Decode(dev.salt)
	key := sa.Add(dev.key)

	pwdPoint.Evaluate(key)
	return pwdPoint.Encode()
}

func TestOPRF(t *testing.T) {
	pwd := "password"
	ctx := &engine.AccountInfo{
		Domain:   "localhost",
		Username: "tester",
	}
	key, _ := newRistrettoKey().Random()
	dev, _ := NewDeviceForAccount(key.Encode(), nil)

	bfac, chal, err := Challenge(pwd, ctx)
	if err != nil {
		t.Error(err)
	}
	resp, err := Response(dev, chal)
	if err != nil {
		t.Error(err)
	}
	rwd, err := Finish(bfac, resp)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(rwd, prf(pwd, ctx, dev)) != 0 {
		t.Error("invalid evaluation")
	}
}
