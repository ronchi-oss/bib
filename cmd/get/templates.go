package get

import (
	"fmt"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "Get all templates",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		templates, err := conf.GetTemplates(targetDir)
		if err != nil {
			return err
		}
		for _, t := range templates {
			fmt.Println(t)
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(templatesCmd)
	utils.InitTargetDirScopedFlags(templatesCmd, &TargetDir, &ProfileName)
}
