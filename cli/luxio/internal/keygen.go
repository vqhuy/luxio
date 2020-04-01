package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vqhuy/luxio/crypto"
	"github.com/vqhuy/luxio/crypto/format"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a random device key",
	Long:  `Generate a random device key`,
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := crypto.GenerateKey()
		encoder := format.ArmoredWriter(os.Stdout)
		encoder.Write(key)
		encoder.Close()
	},
}

func init() {
	RootCmd.AddCommand(keygenCmd)
}
