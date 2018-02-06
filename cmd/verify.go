package cmd

import (
	"log"

	verifyPkg "github.com/flood-io/cli/cmd/verify"
	"github.com/spf13/cobra"
)

var verify verifyPkg.VerifyCmd

// loginCmd represents the login command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "flood verify allows you to verify and debug your flood scripts during development",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if verify.LaunchDevtoolsMode {
			err = verify.LaunchDevtools()
		} else {
			err = verify.Run(args[0])
		}
		if err != nil {
			log.Fatalf("failed to run verify %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVar(&verify.FloodChromeChannel, "channel", "beta", "launch the latest flood chrome on <channel>. Default: beta")

	// hidden dev-only flags
	verifyCmd.Flags().StringVar(&verify.Host, "host", "https://depth.flood.io", "")
	verifyCmd.Flags().StringVar(&verify.DevMode, "devmode", "", "")

	verifyCmd.Flags().MarkHidden("devmode")
	verifyCmd.Flags().MarkHidden("host")
}
