package edit

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note [flags] <note-id>",
	Short: "Edit a note Markdown contents",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		noteID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("can't convert <note-id> argument '%s' to integer", args[0])
		}
		e, err := utils.GetPreferredEditor()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return nil
		}
		notePath := fmt.Sprintf("%s/src/%d/README.md", targetDir, noteID)
		c := exec.Command(e, notePath)
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
	editCmd.AddCommand(noteCmd)
	utils.InitTargetDirScopedFlags(noteCmd, &TargetDir, &ProfileName)
}
