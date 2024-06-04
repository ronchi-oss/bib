package create

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/db"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Create a template",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		if err := db.CreateTemplate(targetDir, args[0]); err != nil {
			return fmt.Errorf("failed creating template file: %v", err)
		}
		fmt.Printf("Template '%s' created successfully\n", args[0])
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
	createCmd.AddCommand(templateCmd)
	utils.InitTargetDirScopedFlags(templateCmd, &TargetDir, &ProfileName)
}
