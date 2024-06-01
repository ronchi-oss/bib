package show

import (
	"github.com/ronchi-oss/bib/cmd"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a single entity",
	Long:  ``,
}

var (
	ProfileName string
	TargetDir   string
)

func init() {
	cmd.RootCmd.AddCommand(showCmd)
}
