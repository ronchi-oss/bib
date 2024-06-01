package create

import (
	"github.com/ronchi-oss/bib/cmd"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an entity",
	Long:  ``,
}

var (
	ProfileName string
	TargetDir   string
)

func init() {
	cmd.RootCmd.AddCommand(createCmd)
}
