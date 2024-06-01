package show

import (
	"fmt"

	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile <name>",
	Short: "Show a single target profile definition",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := conf.GetProfile(args[0])
		if err != nil {
			return fmt.Errorf("failed loading profile '%s': %v", args[0], err)
		}
		d, err := conf.YAMLEncode(profile)
		if err != nil {
			return fmt.Errorf("failed encoding profile to YAML: %v", err)
		}
		fmt.Print(string(d))
		return nil
	},
}

func init() {
	showCmd.AddCommand(profileCmd)
}
