package cmd

import (
	"os"

	"github.com/jplanckeel/kit/pkg/config"
	"github.com/spf13/cobra"
)

var flags config.Flags

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kit",
	Short: "(K)loud install tools",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flags.GithubToken, "github-token", "", "Personal Access Token Github")
}
