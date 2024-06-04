package db

import (
	"fmt"
	"os"
)

func CreateTemplate(targetDir, name string) error {
	path := fmt.Sprintf("%s/tpl/%s", targetDir, name)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("template file '%s' already exists", path)
	}
	if _, err := os.Create(path); err != nil {
		return fmt.Errorf("failed creating template file '%s': %v", path, err)
	}
	return nil
}
