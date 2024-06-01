package generate

import (
	"fmt"

	"github.com/ronchi-oss/bib/conf"
	"github.com/spf13/cobra"
)

var defaultGlobalConfigCmd = &cobra.Command{
	Use:   "default-global-config",
	Short: "Print the default global config file contents",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := conf.GlobalConf{}
		d, err := conf.YAMLEncode(&c)
		if err != nil {
			return fmt.Errorf("failed encoding global config to YAML: %v", err)
		}
		fmt.Println(string(d))
		return nil
	},
}

func init() {
	generateCmd.AddCommand(defaultGlobalConfigCmd)
}
