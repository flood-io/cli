package cmd

import (
	"log"

	initPkg "github.com/flood-io/cli/cmd/init"
	"github.com/spf13/cobra"
)

var initImpl initPkg.InitCmd

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises a Flood Chrome test script and TypeScript environment",
	Long:  `TODO`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		err := initImpl.Run(name)
		if err != nil {
			log.Fatalf("failed to run flood init: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// verifyCmd.Flags().StringVar(&verify.FloodChromeChannel, "channel", "beta", "launch the latest flood chrome on <channel>.")
	// verifyCmd.Flags().BoolVarP(&verify.Verbose, "verbose", "v", false, "print a lot of messages")

	// // hidden dev-only flags
	// verifyCmd.Flags().StringVar(&verify.Host, "host", "https://depth.flood.io", "")
	// verifyCmd.Flags().StringVar(&verify.DevMode, "devmode", "", "")

	// verifyCmd.Flags().MarkHidden("devmode")
	// verifyCmd.Flags().MarkHidden("host")
}
