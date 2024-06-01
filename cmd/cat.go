package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/spf13/cobra"
)

var catCmd = &cobra.Command{
	Use:   "cat [flags] <note-id>",
	Short: "Print a note Markdown contents",
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
		notePath := fmt.Sprintf("%s/src/%d/README.md", targetDir, noteID)
		data, err := os.ReadFile(notePath)
		if err != nil {
			return fmt.Errorf("failed reading README.md file: %v", err)
		}
		fmt.Print(string(data))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(catCmd)
	utils.InitTargetDirScopedFlags(catCmd, &TargetDir, &ProfileName)
}
