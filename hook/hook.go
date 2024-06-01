package hook

import (
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/ronchi-oss/bib/cmd/utils"
	"github.com/ronchi-oss/bib/conf"
)

func NotifyAll(hooks []*conf.Hook, eventName string, args []string) error {
	for _, h := range hooks {
		if !slices.Contains(h.Events, eventName) {
			continue
		}
		path, err := utils.ExpandPath(h.Cmd)
		if err != nil {
			return fmt.Errorf("failed to expand path %s: %v", h.Cmd, err)
		}
		cmd := exec.Command(path, slices.Concat([]string{eventName}, args)...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command '%s' failed: %v", path, err)
		}
	}
	return nil
}
