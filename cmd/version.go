package cmd

import (
	"fmt"

	au "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var Version string
var Commit string
var Date string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the version of flood cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(au.Blue("~~~ Flood Chrome ~~~"))
		fmt.Println(au.Gray("version :"), Version)
		fmt.Println(au.Gray("commit  :"), Commit)
		fmt.Println(au.Gray("date    :"), Date)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
