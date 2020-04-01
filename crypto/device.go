package crypto

type deviceForAccount struct {
	key  *ristrettoKey
	salt []byte
}

// NewDeviceForAccount converts the key material and salt of an account to
// a proper cryptograhic format.
func NewDeviceForAccount(key []byte, salt []byte) (*deviceForAccount, error) {
	devkey := newRistrettoKey()
	if err := devkey.Decode(key); err != nil {
		return nil, err
	}
	if salt != nil {
		return &deviceForAccount{key: devkey, salt: salt}, nil
	}
	devsalt, err := GenerateSalt()
	if err != nil {
		return nil, err
	}
	return &deviceForAccount{key: devkey, salt: devsalt}, nil
}

func generateKey() ([]byte, error) {
	key, err := newRistrettoKey().Random()
	if err != nil {
		return nil, err
	}
	return key.Encode(), nil
}

// GenerateSalt returns a new random salt.
func GenerateSalt() ([]byte, error) {
	return generateKey()
}

// GenerateKey generates a new random device key.
func GenerateKey() ([]byte, error) {
	return generateKey()
}
