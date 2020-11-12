package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "0.1.1"
)

var verCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of luxio",
	Long:  `Print the version number of luxio`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("luxio version " + version)
	},
}

func init() {
	RootCmd.AddCommand(verCmd)
}
