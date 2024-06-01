package get

import (
	"fmt"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

var filtersCmd = &cobra.Command{
	Use:   "filters",
	Short: "Get all filters",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		c, err := conf.GetLocalConf(targetDir)
		if err != nil {
			return err
		}
		for _, filter := range c.Filters {
			fmt.Println(filter.Name)
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(filtersCmd)
	utils.InitTargetDirScopedFlags(filtersCmd, &TargetDir, &ProfileName)
}
