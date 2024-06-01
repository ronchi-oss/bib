package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "bib",
	Short:        "bib(liothek) for collecting and retrieving notes about anything.",
	Long:         ``,
	SilenceUsage: true,
}

var (
	ProfileName string
	TargetDir   string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
