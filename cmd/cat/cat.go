package cat

import (
	"github.com/ronchi-oss/bib/cmd"
	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "Print a file",
	Long:  ``,
}

var (
	ProfileName string
	TargetDir   string
)

func init() {
	cmd.RootCmd.AddCommand(catCmd)
}
