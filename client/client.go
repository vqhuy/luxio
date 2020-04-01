package client

import (
	"github.com/vqhuy/luxio/crypto"
	"github.com/vqhuy/luxio/engine"
)

var outputFormat = []int{
	engine.PIN,
	engine.PlainLowerCase,
	engine.WithSpecialCharacters}

// Start initializes the interaction with the server to a pseudo-random string
// which will be used later to generate a human-readable password.
func Start(pwd string, domain, username string) ([]byte, []byte, error) {
	acc := &engine.AccountInfo{
		Domain:   domain,
		Username: username,
	}
	return crypto.Challenge(pwd, acc)
}

// Finish returns one (or more) human-readable password(s),
// depending on `pwdFormat` (bit-masked).
func Finish(pwdFormat int, bfac []byte, resp []byte) ([]string, error) {
	rwp, err := crypto.Finish(bfac, resp)
	if err != nil {
		return nil, err
	}

	var pwds []string
	for i, f := range outputFormat {
		if ((1 << i) & pwdFormat) == 0 {
			continue
		}
		pwd, err := engine.Generate(rwp, f)
		if err != nil {
			return nil, err
		}
		pwds = append(pwds, pwd)
	}
	return pwds, nil
}
