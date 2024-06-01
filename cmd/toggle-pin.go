package cmd

import (
	"fmt"
	"strconv"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/db"
	"github.com/spf13/cobra"
)

var togglePinCmd = &cobra.Command{
	Use:   "toggle-pin <note-id>",
	Short: "Toggle pin on a note",
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
		return db.TogglePin(targetDir, noteID)
	},
}

func init() {
	RootCmd.AddCommand(togglePinCmd)
	utils.InitTargetDirScopedFlags(togglePinCmd, &TargetDir, &ProfileName)
}
