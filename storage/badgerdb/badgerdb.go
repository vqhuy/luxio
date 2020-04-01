package badgerdb

import (
	"bytes"
	"strconv"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/vqhuy/luxio/engine"
	"github.com/vqhuy/luxio/storage"
)

type badgerdb struct {
	*badger.DB
	keyFunc      func(*engine.AccountInfo, []byte) []byte
	hideMetadata bool
}

var _ storage.DB = (*badgerdb)(nil)

func readKeyPolicy(db *badger.DB, hide bool) (bool, error) {
	// read the policy from the database
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(storage.KeyPolicy))
		switch {
		case err == badger.ErrKeyNotFound:
			// it's the first time, write to the database.
			return txn.Set([]byte(storage.KeyPolicy), []byte(strconv.FormatBool(hide)))
		case err != badger.ErrKeyNotFound && err != nil:
			return err
		default:
			itemCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			cur, err := strconv.ParseBool(string(itemCopy))
			if err != nil {
				return err
			}
			hide = cur // use the existed one instead
		}
		return nil
	})
	return hide, err
}

// OpenDB opens a database from path, with `hide` indicates whether the
// database should hide user's metadata.
// This flag is only used the first time the database is created.
func OpenDB(path string, hide bool) (storage.DB, error) {
	options := badger.DefaultOptions(path)
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	hide, err = readKeyPolicy(db, hide)
	if err != nil {
		return nil, err
	}

	if hide {
		return &badgerdb{
			DB:           db,
			hideMetadata: hide,
			keyFunc:      storage.ComputeKeyHashed,
		}, nil
	}
	return &badgerdb{
		DB:           db,
		hideMetadata: hide,
		keyFunc:      storage.ComputeKey,
	}, nil
}

func (db *badgerdb) Get(key *engine.AccountInfo, ad []byte) ([]byte, error) {
	var salt []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(db.keyFunc(key, ad))
		if err != nil {
			return err
		}
		salt, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func (db *badgerdb) Put(key *engine.AccountInfo, ad []byte, value []byte) error {
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(db.keyFunc(key, ad)), value)
		return err
	})
	return err
}

func (db *badgerdb) Iterate() ([][]byte, error) {
	if db.hideMetadata {
		return nil, storage.ErrUnsupportedOperation
	}
	var list [][]byte
	keyPolicy := []byte(storage.KeyPolicy)
	db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item().KeyCopy(nil)
			if !bytes.Equal(item, keyPolicy) {
				list = append(list, item)
			}
		}
		return nil
	})
	return list, nil
}

func (db *badgerdb) Close() error {
	return db.DB.Close()
}
