package storage

import (
	"errors"

	"github.com/vqhuy/luxio/engine"
	"golang.org/x/crypto/sha3"
)

// hopefully `|` never appear in a domain or an username.
const sep = "|"

// A database should contain a key/value pair which decides if
// the key of each record will be hashed or not.
const KeyPolicy = "key_policy"

var ErrUnsupportedOperation = errors.New("The database is set to hide all the metadata. This operation is unsupported.")

type DB interface {
	Get(key *engine.AccountInfo, ad []byte) ([]byte, error)
	Put(key *engine.AccountInfo, ad []byte, value []byte) error
	Iterate() ([][]byte, error)
	Close() error
}

// ComputeKeyHashed computes SHAKE128 hash of the account
// (domain and username),
// along with some additional data `ad` to hide metadata from the database
// (which could be made public).
// In this implementation, `ad` is the device's key metarial.
// Since the device's key is a cryptographic key, the key stored in
// the database should have enough entropy to prevent brute-force attacks.
// If the dabatase is stored this way, the user cannot query his/her account
// information and might not be able to retrieve his/her password if
// the domain and/or the username is incorrect.
func ComputeKeyHashed(ctx *engine.AccountInfo, ad []byte) []byte {
	h := sha3.NewShake128()
	h.Write(ad)
	h.Write([]byte(ctx.Domain + sep + ctx.Username))
	ret := make([]byte, 32)
	h.Read(ret)
	return ret
}

// ComputeKey generates a database's key from the account info
// (domain and username).
// Anyone has the database can iterate over it and get all the metadata
// including list of domains and usernames.
// However, it might provide more usability since the user is now able to
// query his/her account information.
// The second parameter will be ignored.
func ComputeKey(ctx *engine.AccountInfo, ad []byte) []byte {
	return ctx.Serialize()
}
