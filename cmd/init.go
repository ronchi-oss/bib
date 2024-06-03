package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ronchi-oss/bib/conf"
	"github.com/ronchi-oss/bib/db"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <target-directory>",
	Short: "Initialize <target-directory> as a bib project",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		var targetDir string
		switch {
		case args[0] == ".":
			targetDir = os.Getenv("PWD")
		default:
			targetDir = args[0]
		}

		if _, err := os.Stat(targetDir); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("%s is not a directory.", targetDir)
			}
		}

	loop:
		for {
			fmt.Printf("Will init a new bib project in %s. Continue? (yes/n): ", targetDir)
			reader := bufio.NewReader(os.Stdin)
			answer, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("can't read from standard input.")
			}

			switch strings.TrimSuffix(answer, "\n") {
			case "n":
				fmt.Fprintln(os.Stderr, "Ok, exiting")
				return nil
			case "yes":
				break loop
			}
		}

		fmt.Println("Ok, let's do this.")

		c := conf.LocalConf{
			DefaultTemplate: "default.tpl.md",
			Filters: []*conf.Filter{
				{Name: "all", Cmd: "cat", CmdArgs: []string{}},
				{Name: "pinned", Cmd: "bib-pin-filter", CmdArgs: []string{}},
			},
			Hooks: []*conf.Hook{},
		}

		if err := os.MkdirAll(fmt.Sprintf("%s/tpl", targetDir), 0755); err != nil {
			return fmt.Errorf("failed creating templates directory: %v", err)
		}

		localConfPath := fmt.Sprintf("%s/bib.yml", targetDir)
		defaultTplPath := fmt.Sprintf("%s/tpl/%s", targetDir, c.DefaultTemplate)
		if _, err := os.OpenFile(localConfPath, os.O_RDONLY|os.O_CREATE|os.O_EXCL, 0644); err != nil {
			return fmt.Errorf("target directory config file %s already exists.", localConfPath)
		}
		if _, err := os.OpenFile(defaultTplPath, os.O_RDONLY|os.O_CREATE|os.O_EXCL, 0644); err != nil {
			return fmt.Errorf("target directory default template file %s already exists.", defaultTplPath)
		}

		d, err := conf.YAMLEncode(&c)
		if err != nil {
			return fmt.Errorf("failed encoding local config file to YAML: %v", err)
		}
		if err := ioutil.WriteFile(localConfPath, d, 0644); err != nil {
			return fmt.Errorf("failed writing to local config file: %v", err)
		}

		tpl := []byte(conf.GenerateDefaultTemplateContents())
		if err := ioutil.WriteFile(defaultTplPath, tpl, 0644); err != nil {
			return fmt.Errorf("failed writing to default template file: %v", err)
		}

		if _, err := db.AppendNote(targetDir, "My first note"); err != nil {
			return fmt.Errorf("failed creating first note: %s", err)
		}

		oldwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed determining working directory.")
		}
		if err := os.Chdir(targetDir); err != nil {
			return fmt.Errorf("could not change into %s: %v", targetDir, err)
		}

		if _, err := exec.LookPath("git"); err == nil {
			if _, err := os.Stat(fmt.Sprintf("%s/.git", targetDir)); err != nil && os.IsNotExist(err) {
				commands := []*exec.Cmd{
					exec.Command("git", "init"),
					exec.Command("git", "add", "bib.yml"),
					exec.Command("git", "add", "src/"),
					exec.Command("git", "commit", "-m", "First commit"),
				}
				for _, cmd := range commands {
					fmt.Println(cmd)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("command failed: %v", err)
					}
				}
			}
		}

		if err := os.Chdir(oldwd); err != nil {
			return fmt.Errorf("could not change back into original working directory.")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
