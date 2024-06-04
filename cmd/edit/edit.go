package edit

import (
	"github.com/ronchi-oss/bib/cmd"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a file using BIB_EDITOR (default) or EDITOR (fallback)",
	Long:  ``,
}

var (
	ProfileName string
	TargetDir   string
)

func init() {
	cmd.RootCmd.AddCommand(editCmd)
}
