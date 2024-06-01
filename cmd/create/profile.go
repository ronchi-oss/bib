package create

import (
	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile <name> <target-directory>",
	Short: "Create a profile",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(2)),
	RunE: func(cmd *cobra.Command, args []string) error {
		return conf.AppendProfile(args[0], args[1])
	},
}

func init() {
	createCmd.AddCommand(profileCmd)
}
