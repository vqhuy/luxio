package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/vqhuy/luxio/engine"
	"github.com/vqhuy/luxio/server"
	"github.com/vqhuy/luxio/storage/badgerdb"
)

const (
	hasChildPrefix     = "├──"
	innerPrefix        = "|  "
	lastPrefix         = "└──"
	childPaddingPrefix = "   "
	rootPrefix         = "o"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all account for the given domain",
	Long:  `List all account for the given domain (* for all)`,
	Args:  cobra.ExactArgs(1),
	Run:   list,
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) {
	domain := engine.DomainFromUrl(args[0])

	conf := readConfig()
	db, err := badgerdb.OpenDB(conf.DB, conf.HideMetadata)
	if err != nil {
		log.Fatal(err)
	}

	server, err := server.NewServer(db, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	results, err := server.CmdList(domain)
	if err != nil {
		log.Fatal(err)
	}
	counter := 0
	fmt.Println(rootPrefix)
	for key, value := range results {
		inner := innerPrefix
		childCounter := 0
		if counter == len(results)-1 {
			inner = childPaddingPrefix
			fmt.Println(lastPrefix + key)
		} else {
			fmt.Println(hasChildPrefix + key)
		}
		for index, site := range value {
			if index == len(value)-1 {
				fmt.Println(inner + lastPrefix + site)
			} else {
				fmt.Println(inner + hasChildPrefix + site)
			}
			childCounter++
		}
		counter++
	}
}
