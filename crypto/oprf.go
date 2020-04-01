package crypto

import (
	"crypto/sha512"

	"github.com/vqhuy/luxio/engine"
	"golang.org/x/crypto/hkdf"
)

// Challenge takes the master password `pwd`,
// the account information `ctx` (including username and url)
// and returns a blinding factor and a challenge.
func Challenge(pwd string, ctx *engine.AccountInfo) (
	[]byte, []byte, error) {
	chal := newRistrettoPoint()
	bfac := newRistrettoKey()
	h := hkdf.New(sha512.New, []byte(pwd), nil, ctx.Serialize())
	b := make([]byte, 64)
	if _, err := h.Read(b); err != nil {
		return nil, nil, err
	}
	chal.FromUniformBytes(b)
	bfac.Random()
	chal.Evaluate(bfac)
	return bfac.Encode(), chal.Encode(), nil
}

// Response takes the device's key (and salt) `dev`,
// the challenge `chal` from the client,
// and returns a response.
func Response(dev *deviceForAccount, chal []byte) ([]byte, error) {
	ch := newRistrettoPoint()
	if err := ch.Decode(chal); err != nil {
		return nil, err
	}
	sa := newRistrettoKey()
	if err := sa.Decode(dev.salt); err != nil {
		return nil, err
	}
	p := sa.Add(dev.key)
	return ch.Evaluate(p).Encode(), nil
}

// Finish takes the blinding factor generated from `Challenge` and
// the response from `Response` and output a pseudo-random string.
func Finish(bfac []byte, resp []byte) ([]byte, error) {
	bf := newRistrettoKey()
	if err := bf.Decode(bfac); err != nil {
		return nil, err
	}
	re := newRistrettoPoint()
	if err := re.Decode(resp); err != nil {
		return nil, err
	}
	bf.Invert()
	rwd := re.Evaluate(bf)
	return rwd.Encode(), nil
}
