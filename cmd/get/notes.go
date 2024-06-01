package get

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/conf"
	"github.com/ronchi-oss/bib/db"
	"github.com/spf13/cobra"
)

var filterName string

var notesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Get all notes of a bib target directory",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		notes, err := db.GetNotes(targetDir)
		if err != nil {
			return fmt.Errorf("failed loading notes: %v", err)
		}
		var filterCmd *exec.Cmd
		reader, writer := os.Stdin, os.Stdout
		if len(filterName) > 0 {
			f, err := conf.GetFilter(targetDir, filterName)
			if err != nil {
				return fmt.Errorf("failed loading filter '%s': %v", filterName, err)
			}
			reader, writer, err = os.Pipe()
			if err != nil {
				return fmt.Errorf("unexpected error: %v", err)
			}
			path, err := utils.ExpandPath(f.Cmd)
			if err != nil {
				return fmt.Errorf("failed to expand path %s: %v", f.Cmd, err)
			}
			filterCmd = exec.Command(path, f.CmdArgs...)
			filterCmd.Stdin = reader
			filterCmd.Stdout = os.Stdout
			filterCmd.Stderr = os.Stderr
		}
		for _, note := range notes {
			pinned := " "
			if note.Pinned {
				pinned = "*"
			}
			fmt.Fprintf(writer, "%d\t%s\t%s\t%s\n",
				note.ID,
				pinned,
				note.CreatedAt.Format(time.DateTime),
				note.Title)
		}
		writer.Close()
		if filterCmd != nil {
			if err := filterCmd.Run(); err != nil {
				return fmt.Errorf("filter process '%s' failed: %v", err)
			}
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(notesCmd)
	utils.InitTargetDirScopedFlags(notesCmd, &TargetDir, &ProfileName)
	notesCmd.Flags().StringVarP(&filterName, "filter", "f", "", "Result set filter")
	notesCmd.RegisterFlagCompletionFunc("filter", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		td, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		c, err := conf.GetLocalConf(td)
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		fNames := []string{}
		for _, f := range c.Filters {
			fNames = append(fNames, f.Name)
		}
		return fNames, cobra.ShellCompDirectiveDefault
	})
}
