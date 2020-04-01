package server

import (
	"os/exec"
	"strings"

	"github.com/vqhuy/luxio/crypto"
	"github.com/vqhuy/luxio/engine"
	"github.com/vqhuy/luxio/storage"
)

type server struct {
	db  storage.DB
	key []byte
}

// NewServer creates a new server instance from the database `db`
// and the device key `key`.
func NewServer(db storage.DB, key []byte) (*server, error) {
	s := new(server)
	s.db = db
	s.key = key
	return s, nil
}

func (s *server) Shutdown() error {
	return s.db.Close()
}

func (s *server) CmdRequest(domain, username string, chal []byte) ([]byte, error) {
	acc := &engine.AccountInfo{
		Domain:   domain,
		Username: username,
	}
	// get salt from database
	salt, err := s.db.Get(acc, s.key)
	if err != nil {
		return nil, err
	}
	dev, err := crypto.NewDeviceForAccount(s.key, salt)
	if err != nil {
		return nil, err
	}

	return crypto.Response(dev, chal)
}

func (s *server) CmdUpdate(domain, username string, chal []byte) ([]byte, error) {
	acc := &engine.AccountInfo{
		Domain:   domain,
		Username: username,
	}
	// generate new salt
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, err
	}
	dev, err := crypto.NewDeviceForAccount(s.key, salt)
	if err != nil {
		return nil, err
	}

	// get response
	resp, err := crypto.Response(dev, chal)
	if err != nil {
		return nil, err
	}

	// and only write to db if everything is okay so far
	if err := s.db.Put(acc, s.key, salt); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *server) CmdList(search string) (map[string][]string, error) {
	list := make(map[string][]string)
	all, err := s.db.Iterate()
	if err != nil {
		return nil, err
	}

	tmp := &engine.AccountInfo{}
	for _, item := range all {
		domain, username := tmp.Deserialize(item)
		if search == "*" || strings.Contains(domain, search) {
			list[domain] = append(list[domain], username)
		}
	}
	return list, nil
}

func executeCommand(cmd string) ([]byte, error) {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyEval evaluates the given command `cmd` and returns the key material.
func KeyEval(cmd string) ([]byte, error) {
	if key, err := executeCommand(cmd); err != nil {
		return nil, err
	} else {
		return []byte(string(key)), nil
	}
}
