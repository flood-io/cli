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
		var err error
		if bluDev.LaunchDevtoolsMode {
			err = bluDev.LaunchDevtools()
		} else {
			err = bluDev.Run(args[0])
		}
		if err != nil {
			log.Fatalf("failed to run dev-blu %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(bluDevCmd)

	bluDevCmd.Flags().StringVar(&bluDev.FloodChromeChannel, "channel", "beta", "launch the latest flood chrome on <channel>. Default: beta")

	// bluDevCmd.Flags().BoolVarP(&bluDev.LaunchDevtoolsMode, "devtools", "d", false, "launch chrome instance devtools connected to flood-chrome BLU dev mode")
}
