package main

import (
	"fmt"
	"os"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/flood-io/cli/cmd"
)

var version string = "dev"
var commit string = "dev"
var date string = "today"
var bugsnagAPIKey string

func init() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date
}

func main() {
	// logger := log.New(os.Stderr, "", log.Lshortfile)

	releaseStage := "development"
	if version != "dev" {
		releaseStage = "production"
	}

	bugsnag.Configure(bugsnag.Configuration{
		// Your Bugsnag project API key
		APIKey: bugsnagAPIKey,

		// The development stage of your application build, like "alpha" or
		// "production"
		ReleaseStage: releaseStage,

		// only notify once released
		NotifyReleaseStages: []string{"production"},

		// The import paths for the Go packages containing your source files
		ProjectPackages: []string{"main", "github.com/flood-io/cli*", "github.com/flood-io/go-wrenches*"},

		Logger: nil,
	})

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("An internal error occurred, sorry!")
			fmt.Println("We've been notified of the problem.")
			fmt.Println("Feel free to try again soon, or if your problem persists please")
			fmt.Println("contact Flood support at support@flood.io")

			if os.Getenv("FLOOD_DEBUG") != "" {
				fmt.Println("flood debug mode on, re-panicing:")
				panic(r)
			} else {
				os.Exit(1)
			}
		}
	}()
	defer bugsnag.AutoNotify()

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
