package cli

import (
	"log"

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
	return &conf
}
