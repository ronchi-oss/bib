package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/spf13/cobra"
)

var title string

var renderCmd = &cobra.Command{
	Use:   "render [flags] <template-name>",
	Short: "Render a template",
	Long: `render renders a template without adding a note to the target directory.

If the template defines a {{ .Body }}, the command reads it from standard input when "-" is provided as the last argument.
`,
	Args: cobra.MatchAll(cobra.RangeArgs(1, 2)),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return utils.TemplateNameShellComp(TargetDir, ProfileName)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		targetDir, err := utils.GetTargetDir(TargetDir, ProfileName)
		if err != nil {
			return fmt.Errorf("failed determining target directory: %v", err)
		}
		var body []byte
		if args[len(args)-1] == "-" {
			body, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed reading standard input: %v", err)
			}
		}
		path := fmt.Sprintf("%s/tpl/%s", targetDir, args[0])
		if err := utils.RenderTemplate(os.Stdout, path, title, body); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(renderCmd)
	utils.InitTargetDirScopedFlags(renderCmd, &TargetDir, &ProfileName)
	renderCmd.Flags().StringVarP(&title, "title", "t", "", "Template {{ .Title }}")
}
