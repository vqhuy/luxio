package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vqhuy/luxio/cli"
	"github.com/vqhuy/luxio/client"
	"github.com/vqhuy/luxio/crypto/format"
	"github.com/vqhuy/luxio/engine"
	"github.com/vqhuy/luxio/server"
	"github.com/vqhuy/luxio/storage/badgerdb"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	configFileName = ".luxiorc"
	luxioDirName   = ".luxio"
	dbDirName      = "db"
	keyFileName    = "key.luxio"
)

func makeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("domain", "d", "", "Domain name (an URL, a website, etc)")
	cmd.MarkFlagRequired("domain")
	cmd.Flags().StringP("username", "u", "", "Username or Account")
	cmd.MarkFlagRequired("username")
	cmd.Flags().Bool("pin", false, "Print as a PIN code")
	cmd.Flags().Bool("plain", false, "Print as a plain, lower-case passphrase")
	cmd.Flags().Bool("special", false, "Print as a title-case passphrase with a fixed postfix")
	cmd.Flags().Bool("password", false, "Print as a random password with a fixed postfix")
}

func luxioHomeDir() string {
	// default home directory
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return homedir
}

func readConfig() *cli.ServerConfig {
	return cli.ReadServerConfig(path.Join(luxioHomeDir(), configFileName))
}

func readDeviceKey(conf *cli.ServerConfig) []byte {
	var reader *bytes.Reader
	switch {
	case conf.Key != "" && conf.KeyEval != "":
		log.Fatal("Configuration has both Key and KeyEval. Choose only one.")
	case conf.Key != "":
		key := []byte(conf.Key)
		reader = bytes.NewReader(key)
	case conf.KeyEval != "":
		key, err := server.KeyEval(conf.KeyEval)
		if err != nil {
			log.Fatal(err)
		}
		reader = bytes.NewReader(key)
	default:
		log.Fatal("No Device Key.")
	}
	decoder := format.ArmoredReader(reader)
	key, _ := ioutil.ReadAll(decoder)
	return key
}

func handle(cmd *cobra.Command, args []string) {
	pwd, err := readMasterPwd()
	if err != nil {
		log.Fatal(err)
	}
	domain := engine.DomainFromUrl(cmd.Flag("domain").Value.String())
	username := cmd.Flag("username").Value.String()

	pin, err := strconv.ParseBool(cmd.Flag("pin").Value.String())
	if err != nil {
		log.Fatal(err)
	}
	plain, err := strconv.ParseBool(cmd.Flag("plain").Value.String())
	if err != nil {
		log.Fatal(err)
	}
	special, err := strconv.ParseBool(cmd.Flag("special").Value.String())
	if err != nil {
		log.Fatal(err)
	}
	password, err := strconv.ParseBool(cmd.Flag("password").Value.String())
	if err != nil {
		log.Fatal(err)
	}
	if !pin && !plain && !special && !password {
		plain = true
	}

	conf := readConfig()
	db, err := badgerdb.OpenDB(conf.DB, conf.HideMetadata)
	if err != nil {
		log.Fatal(err)
	}

	server, err := server.NewServer(db, readDeviceKey(conf))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	bfac, chal, err := client.Start(pwd, domain, username)
	if err != nil {
		log.Fatal(err)
	}

	var resp []byte
	if cmd.Name() == "update" {
		resp, err = server.CmdUpdate(domain, username, chal)
		if err != nil {
			log.Fatal(err)
		}
	} else if cmd.Name() == "request" {
		resp, err = server.CmdRequest(domain, username, chal)
		if err != nil {
			log.Fatal(err)
		}
	}
	format := 0
	if pin {
		format ^= 1 << engine.PIN
	}
	if plain {
		format ^= 1 << engine.PlainLowerCase
	}
	if special {
		format ^= 1 << engine.WithSpecialCharacters
	}
	if password {
		format ^= 1 << engine.Password
	}
	rwd, err := client.Finish(format, bfac, resp)
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range rwd {
		fmt.Println(w)
	}
}

func readMasterPwd() (string, error) {
	fd := int(os.Stdin.Fd())
	state, err := terminal.GetState(fd)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(fd, state)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		terminal.Restore(fd, state)
		os.Exit(1)
	}()

	fmt.Print("❯ Enter your Master Password:")
	pwd, err := terminal.ReadPassword(fd)
	fmt.Println("")
	return string(pwd), err
}
