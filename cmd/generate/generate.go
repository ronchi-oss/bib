package generate

import (
	"github.com/ronchi-oss/bib/cmd"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate specific configuration",
	Long:  ``,
}

func init() {
	cmd.RootCmd.AddCommand(generateCmd)
}
