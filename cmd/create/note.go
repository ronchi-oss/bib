package create

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/conf"
	"github.com/ronchi-oss/bib/db"
	"github.com/ronchi-oss/bib/hook"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Create a note",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		title := "My new note"
		id, err := db.AppendNote(targetDir, title)
		if err != nil {
			return fmt.Errorf("failed creating note files: %v", err)
		}
		fmt.Printf("Note created successfully (id = %d)\n", id)
		conf, err := conf.GetLocalConf(targetDir)
		if err != nil {
			return fmt.Errorf("failed loading target directory bib config file: %v", err)
		}
		hook.NotifyAll(conf.Hooks, "note.created", []string{strconv.Itoa(id), title})
		e, err := utils.GetPreferredEditor()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return nil
		}
		notePath := fmt.Sprintf("%s/src/%d/README.md", targetDir, id)
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
	createCmd.AddCommand(noteCmd)
	utils.InitTargetDirScopedFlags(noteCmd, &TargetDir, &ProfileName)
}
