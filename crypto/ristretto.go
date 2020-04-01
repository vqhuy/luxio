package crypto

import (
	"crypto/rand"

	"github.com/gtank/ristretto255"
)

type ristrettoKey struct {
	s *ristretto255.Scalar
}

type ristrettoPoint struct {
	e *ristretto255.Element
}

func newRistrettoKey() *ristrettoKey {
	return &ristrettoKey{
		s: ristretto255.NewScalar(),
	}
}

func (key *ristrettoKey) Random() (*ristrettoKey, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return key, err
	}
	key.s.FromUniformBytes(b)
	return key, nil
}

func (key *ristrettoKey) Add(other *ristrettoKey) *ristrettoKey {
	key.s.Add(key.s, other.s)
	return key
}

func (key *ristrettoKey) Subtract(subtrahend *ristrettoKey) *ristrettoKey {
	key.s.Subtract(key.s, subtrahend.s)
	return key
}

func (key *ristrettoKey) Invert() *ristrettoKey {
	key.s.Invert(key.s)
	return key
}

func (key *ristrettoKey) FromUniformBytes(b []byte) *ristrettoKey {
	key.s.FromUniformBytes(b)
	return key
}

func (key *ristrettoKey) Encode() []byte {
	return key.s.Encode(nil)
}

func (key *ristrettoKey) Decode(in []byte) error {
	err := key.s.Decode(in)
	return err
}

func newRistrettoPoint() *ristrettoPoint {
	return &ristrettoPoint{
		e: ristretto255.NewElement(),
	}
}

func (point *ristrettoPoint) Evaluate(key *ristrettoKey) *ristrettoPoint {
	point.e.ScalarMult(key.s, point.e)
	return point
}

func (point *ristrettoPoint) FromUniformBytes(b []byte) *ristrettoPoint {
	point.e.FromUniformBytes(b)
	return point
}

func (point *ristrettoPoint) Encode() []byte {
	return point.e.Encode(nil)
}

func (point *ristrettoPoint) Decode(in []byte) error {
	err := point.e.Decode(in)
	return err
}
