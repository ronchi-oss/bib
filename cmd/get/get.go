package get

import (
	"github.com/ronchi-oss/bib/cmd"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a list of entities",
	Long:  ``,
}

var (
	ProfileName string
	TargetDir   string
)

func init() {
	cmd.RootCmd.AddCommand(getCmd)
}
