package cmd

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Add or update password",
	Long:    `Add or update password for the given account on the given domain`,
	Example: `luxio update -d "https://accounts.google.com/" -u "name@gmail.com"`,
	Run:     handle,
}

func init() {
	makeFlags(updateCmd)
	RootCmd.AddCommand(updateCmd)
}
