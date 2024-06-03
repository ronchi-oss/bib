package cat

import (
	"fmt"
	"os"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template [flags] <template-name>",
	Short: "Print a template file contents",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		path := fmt.Sprintf("%s/tpl/%s", targetDir, args[0])
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed reading template file '%s': %v", path, err)
		}
		fmt.Print(string(data))
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		templates, err := conf.GetTemplates(targetDir)
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		return templates, cobra.ShellCompDirectiveDefault
	},
}

func init() {
	catCmd.AddCommand(templateCmd)
	utils.InitTargetDirScopedFlags(templateCmd, &TargetDir, &ProfileName)
}
