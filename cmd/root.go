package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gourd",
	Short: "Serve interview questions dynamically from private repositories",
	Long: `gourd provides the ability to serve interview questions written in markdown by providing URL and PAT for the
	requested repository. Solutions submitted by users will be commited automatically to a unique branch for review`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
