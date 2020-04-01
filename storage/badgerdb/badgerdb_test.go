package badgerdb

import (
	"io/ioutil"
	"os"
	"testing"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/vqhuy/luxio/engine"
	"github.com/vqhuy/luxio/storage"
)

func setupBadgerDB(t *testing.T, preTest func(t *testing.T, db *badger.DB),
	postTest func(t *testing.T, db *badger.DB)) {
	dir, _ := ioutil.TempDir("", "luxio-badgerdb-test")
	defer func() {
		os.RemoveAll(dir)
	}()

	options := badger.DefaultOptions(dir)
	options.Logger = nil
	db, _ := badger.Open(options)
	preTest(t, db)
	db.Close()
	db, _ = badger.Open(options)
	postTest(t, db)
	db.Close()
}

func TestReadKeyPolicy(t *testing.T) {
	for _, tc := range []struct {
		firstCall bool
		firstWant bool
		nextCall  bool
		nextWant  bool
	}{
		{true, true, true, true},
		{true, true, false, true},
		{false, false, false, false},
		{false, false, true, false},
	} {
		setupBadgerDB(t, func(t *testing.T, db *badger.DB) {
			got, err := readKeyPolicy(db, tc.firstCall)
			if err != nil {
				t.Error(err)
			}
			if got != tc.firstWant {
				t.Error("expect", tc.firstWant, "got", got)
			}
		}, func(t *testing.T, db *badger.DB) {
			got, err := readKeyPolicy(db, tc.nextCall)
			if err != nil {
				t.Error(err)
			}
			if got != tc.nextWant {
				t.Error("expect", tc.nextWant, "got", got)
			}
		})
	}
}

func setupDB(t *testing.T, hide bool, test func(t *testing.T, db *badgerdb)) {
	dir, _ := ioutil.TempDir("", "luxio-badgerdb-test")
	defer func() {
		os.RemoveAll(dir)
	}()
	db, _ := OpenDB(dir, hide)
	test(t, db.(*badgerdb))
	db.Close()
}

func TestDBFollowKeyPolicy(t *testing.T) {
	acc := &engine.AccountInfo{
		Domain:   "test.com",
		Username: "tester",
	}
	ad := []byte{0}

	for _, tc := range []struct {
		hide bool
		key  []byte
		want error
	}{
		{true, storage.ComputeKey(acc, ad), badger.ErrKeyNotFound},
		{true, storage.ComputeKeyHashed(acc, ad), nil},
		{false, storage.ComputeKeyHashed(acc, ad), badger.ErrKeyNotFound},
		{false, storage.ComputeKey(acc, ad), nil},
	} {
		setupDB(t, tc.hide, func(t *testing.T, db *badgerdb) {
			db.Put(acc, ad, []byte{0})

			err := db.View(func(txn *badger.Txn) error {
				_, err := txn.Get(tc.key)
				return err
			})
			if err != tc.want {
				t.Error("expect", tc.want, "got", err)
			}
		})
	}
}

func TestDBIterate(t *testing.T) {
	acc1 := &engine.AccountInfo{
		Domain:   "test.com",
		Username: "tester1",
	}
	acc2 := &engine.AccountInfo{
		Domain:   "test.com",
		Username: "tester2",
	}
	for _, tc := range []struct {
		hide       bool
		wantErr    error
		wantLength int
	}{
		{true, storage.ErrUnsupportedOperation, 0},
		{false, nil, 2},
	} {
		setupDB(t, tc.hide, func(t *testing.T, db *badgerdb) {
			db.Put(acc1, []byte{0}, []byte{0})
			db.Put(acc2, []byte{0}, []byte{0})
			res, err := db.Iterate()
			if err != tc.wantErr {
				t.Error("expect", tc.wantErr, "got", err)
			}
			if res != nil && len(res) != tc.wantLength {
				t.Error("expect", tc.wantLength, "got", len(res))
			}
		})
	}
}
