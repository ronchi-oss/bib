package edit

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template [flags] <template-name>",
	Short: "Edit a template contents",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return utils.TemplateNameShellComp(TargetDir, ProfileName)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		e, err := utils.GetPreferredEditor()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return nil
		}
		tplPath := fmt.Sprintf("%s/tpl/%s", targetDir, args[0])
		c := exec.Command(e, tplPath)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed running editor command '%s': %v", e, err)
		}
		return nil
	},
}

func init() {
	editCmd.AddCommand(templateCmd)
	utils.InitTargetDirScopedFlags(templateCmd, &TargetDir, &ProfileName)
}
