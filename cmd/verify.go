package cmd

import (
	"fmt"
	"os"

	verifyPkg "github.com/flood-io/cli/cmd/verify"
	"github.com/flood-io/cli/config"
	"github.com/flood-io/cli/oauthclient"
	au "github.com/logrusorgru/aurora"
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
		fmt.Println(au.Blue("~~~ Flood Chrome"), au.Green("Verify"), au.Blue("~~~"))
		var err error

		authToken, err := oauthclient.GetAuthToken(config.DefaultAuthCache())
		if err != nil {
			fmt.Printf("Failed to get auth token: %s\n", err)
			fmt.Println("Please log in with", au.Gray("flood login"))
			os.Exit(1)
		}

		if verify.LaunchDevtoolsMode {
			err = verify.LaunchDevtools()
		} else {
			err = verify.Run(authToken, args[0])
		}
		if err != nil {
			fmt.Printf("failed to run verify %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVar(&verify.FloodChromeChannel, "channel", "beta", "launch the latest flood chrome on <channel>.")
	verifyCmd.Flags().BoolVarP(&verify.Verbose, "verbose", "v", false, "print a lot of messages")

	// hidden dev-only flags
	verifyCmd.Flags().StringVar(&verify.Host, "host", "https://depth.flood.io", "")
	verifyCmd.Flags().StringVar(&verify.DevMode, "devmode", "", "")

	verifyCmd.Flags().MarkHidden("devmode")
	verifyCmd.Flags().MarkHidden("host")
}
