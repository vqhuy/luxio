package cmd

import (
	"github.com/spf13/cobra"
)

var requestCmd = &cobra.Command{
	Use:     "request",
	Short:   "Get password",
	Long:    `Get password of the given account on the given domain`,
	Example: `luxio request -d "https://accounts.google.com/" -u "name@gmail.com"`,
	Run:     handle,
}

func init() {
	makeFlags(requestCmd)
	RootCmd.AddCommand(requestCmd)
}
