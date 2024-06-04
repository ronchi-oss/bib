package create

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/conf"
	"github.com/ronchi-oss/bib/db"
	"github.com/ronchi-oss/bib/hook"
	"github.com/spf13/cobra"
)

var (
	template string
	title    string
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Create a note",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err       error
			readStdin bool
			body      []byte
		)
		if len(args) > 0 && args[len(args)-1] == "-" {
			readStdin = true
			body, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed reading standard input: %v", err)
			}
		}
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		c, err := conf.GetLocalConf(targetDir)
		if err != nil {
			return err
		}
		if len(template) == 0 {
			template = c.DefaultTemplate
		}
		id, err := db.AppendNote(targetDir, template, title, body)
		if err != nil {
			return fmt.Errorf("failed creating note files: %v", err)
		}
		fmt.Printf("Note created successfully (id = %d)\n", id)
		conf, err := conf.GetLocalConf(targetDir)
		if err != nil {
			return fmt.Errorf("failed loading target directory bib config file: %v", err)
		}
		hook.NotifyAll(conf.Hooks, "note.created", []string{strconv.Itoa(id), title})
		if !readStdin {
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
		}
		return nil
	},
}

func init() {
	createCmd.AddCommand(noteCmd)
	utils.InitTargetDirScopedFlags(noteCmd, &TargetDir, &ProfileName)
	noteCmd.Flags().StringVarP(&title, "title", "", "", "Note title")
	noteCmd.Flags().StringVarP(&template, "template", "t", "", "Template file (base name) that should render the note")
	noteCmd.RegisterFlagCompletionFunc("template", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return utils.TemplateNameShellComp(TargetDir, ProfileName)
	})
}
