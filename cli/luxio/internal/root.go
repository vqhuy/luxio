package cmd

import cli "github.com/vqhuy/luxio/cli"

// RootCmd represents the base "luxio" command when called without any subcommands.
var RootCmd = cli.NewRootCommand("luxio",
	"Luxio - Another password manager",
	`Luxio - Another password manager`)
