package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print bib version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bib 0.1.0")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
