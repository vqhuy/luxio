package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/vqhuy/luxio/cli"
	"github.com/vqhuy/luxio/crypto"
	"github.com/vqhuy/luxio/crypto/format"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your luxio",
	Long:  `Initialize a default config and a working directory at your home directory. Becareful to not overwrite your current vault!`,
	Run:   mkConfig,
}

func init() {
	RootCmd.AddCommand(initCmd)
}

func mkConfig(cmd *cobra.Command, args []string) {
	dir := path.Join(luxioHomeDir(), luxioDirName)
	dbPath := path.Join(dir, dbDirName)
	keyPath := path.Join(dir, keyFileName)
	if err := os.Mkdir(dir, 0700); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir(dbPath, 0700); err != nil {
		log.Fatal(err)
	}

	// generate a random key
	key, _ := crypto.GenerateKey()
	f, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		log.Fatal(err)
	}
	encoder := format.ArmoredWriter(f)
	encoder.Write(key)
	encoder.Close()

	// generate config file
	conf := &cli.ServerConfig{
		DB:           dbPath,
		KeyEval:      fmt.Sprintf("cat %s", keyPath),
		HideMetadata: false,
	}

	var confBuf bytes.Buffer
	e := toml.NewEncoder(&confBuf)
	if err := e.Encode(conf); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(path.Join(luxioHomeDir(), configFileName),
		confBuf.Bytes(), 0600); err != nil {
		log.Fatal(err)
	}
}
