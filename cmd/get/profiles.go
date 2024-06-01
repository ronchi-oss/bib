package get

import (
	"fmt"

	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Get all target profiles",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		gconf, err := conf.GetGlobalConf()
		if err != nil {
			return fmt.Errorf("failed loading global config: %v", err)
		}
		for _, profile := range gconf.Profiles {
			fmt.Printf("%s\t%s\n", profile.Name, profile.TargetDir)
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(profilesCmd)
}
