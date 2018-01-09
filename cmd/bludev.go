package cmd

import (
	"log"

	"github.com/flood-io/cli/cmd/bludev"
	"github.com/spf13/cobra"
)

var bluDev bludev.BLUDev

// loginCmd represents the login command
var bluDevCmd = &cobra.Command{
	Use:   "dev-blu",
	Short: "Develop & debug yourflood BLU script",
	Run: func(cmd *cobra.Command, args []string) {
		err := bluDev.Run(args[0])
		if err != nil {
			log.Fatalf("failed to run dev-blu %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(bluDevCmd)
}
