package main

import (
	"github.com/vqhuy/luxio/cli"
	cmd "github.com/vqhuy/luxio/cli/luxio/internal"
)

func main() {
	cli.Execute(cmd.RootCmd)
}
