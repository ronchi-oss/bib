package generate

import (
	"fmt"

	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

var defaultTemplateCmd = &cobra.Command{
	Use:   "default-template",
	Short: "Print the default note template contents",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(conf.GenerateDefaultTemplateContents())
		return nil
	},
}

func init() {
	generateCmd.AddCommand(defaultTemplateCmd)
}
