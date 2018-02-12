package cmd

import (
	"fmt"
	"os"

	initPkg "github.com/flood-io/cli/cmd/init"
	au "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var initImpl initPkg.InitCmd

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialises a Flood Chrome test script and TypeScript environment",
	Long: au.Gray("flood init").String() + ` initialises a sample Flood Chrome test
script and TypeScript environment to get you started using Flood Chrome.

The test script 'test.ts' gives you the general structure of a Flood Chrome test
script.

The TypeScript environment provides our recommended configuration for use with 
an intelligent editor such as Microsoft VSCode. Using such an editor to write 
your Flood Chrome test script gives you programming super powers such as 
autocompletion and in-line documentation.

'flood init' won't write files to a non-empty directory unless your force it with
--force.

Example: 
	flood init my-load test --url https://load-test-target.com
	# your test env is initialised

Example: using interactive config
	mkdir my-load-test
	cd my-load-test
	flood init 
	# (now follow the prompts to configure your test scripts:)
	~~~ Flood Chrome Init ~~~
	Project name
	Enter a value (Default is my-load-test): [hit enter]

	Test URL
	Enter a value (Default is https://challenge.flood.io/): [type your test url, hit enter]

	# your test env is initialised

	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		err := initImpl.Run(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to run flood init:", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&initImpl.URL, "url", "https://challenge.flood.io/", "A URL to use in the example template.")
	initCmd.Flags().BoolVar(&initImpl.Force, "force", false, "Force creation even if it may overwrite existing files.")
}
