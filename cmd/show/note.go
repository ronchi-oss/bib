package show

import (
	"fmt"
	"strconv"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/conf"
	"github.com/ronchi-oss/bib/db"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Show a single note's metadata",
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
		note, err := db.GetNote(targetDir, noteID)
		if err != nil {
			return fmt.Errorf("failed loading note: %v", err)
		}
		d, err := conf.YAMLEncode(&note)
		if err != nil {
			return fmt.Errorf("failed encoding note metadata to YAML: %v", err)
		}
		fmt.Println(string(d))
		return nil
	},
}

func init() {
	showCmd.AddCommand(noteCmd)
	utils.InitTargetDirScopedFlags(noteCmd, &TargetDir, &ProfileName)
}
