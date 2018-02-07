package cmd

import (
	"log"

	verifyPkg "github.com/flood-io/cli/cmd/verify"
	"github.com/flood-io/cli/config"
	"github.com/spf13/cobra"
)

var verify verifyPkg.VerifyCmd

// loginCmd represents the login command
var verifyCmd = &cobra.Command{
	Use:   "verify <test-script.ts>",
	Short: "Verifies your Flood Chrome test scripts during development",
	Long: `Use 'flood verify' to verify and debug your Flood Chrome script during development.

Once its ready, use the script to run a full load test at high concurrency via Flood IO.
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authToken := config.DefaultAuthCache().MustToken()

		var err error
		if verify.LaunchDevtoolsMode {
			err = verify.LaunchDevtools()
		} else {
			err = verify.Run(authToken, args[0])
		}
		if err != nil {
			log.Fatalf("failed to run verify %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVar(&verify.FloodChromeChannel, "channel", "beta", "launch the latest flood chrome on <channel>.")

	// hidden dev-only flags
	verifyCmd.Flags().StringVar(&verify.Host, "host", "https://depth.flood.io", "")
	verifyCmd.Flags().StringVar(&verify.DevMode, "devmode", "", "")

	verifyCmd.Flags().MarkHidden("devmode")
	verifyCmd.Flags().MarkHidden("host")
}
