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

	initCmd.Flags().StringVar(&initImpl.URL, "url", "https://challenge.flood.io/", "A URL to use in the example template.")
	initCmd.Flags().BoolVar(&initImpl.Force, "force", false, "Force creation even if it may overwrite existing files.")
}
