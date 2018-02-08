package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "flood",
	Short: "Flood Command Line Interface",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
