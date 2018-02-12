package main

import (
	"fmt"
	"os"

	"github.com/flood-io/cli/cmd"
)

var version string = "dev"
var commit string = "dev"
var date string = "today"

func init() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
