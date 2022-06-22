package cli

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	DB           string
	HideMetadata bool
	Key          string
	KeyEval      string
}

func ReadServerConfig(path string) *ServerConfig {
	var conf ServerConfig
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(conf.DB, "~/") {
		dir, _ := os.UserHomeDir()
		conf.DB = filepath.Join(dir, conf.DB[2:])
	}
	return &conf
}
