package cmd

import (
	"github.com/flood-io/cli/config"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logs out of Flood",
	Run: func(cmd *cobra.Command, args []string) {
		cache := config.DefaultAuthCache()
		cache.Clear()
	},
}

func init() {
	RootCmd.AddCommand(logoutCmd)
}
