package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type LocalConf struct {
	Filters []*Filter `yaml:"filters"`
	Hooks   []*Hook   `yaml:"hooks"`
}

type Filter struct {
	Name    string   `yaml:"name"`
	Cmd     string   `yaml:"cmd"`
	CmdArgs []string `yaml:"cmd_args"`
}

type Hook struct {
	Events []string `yaml:"events"`
	Cmd    string   `yaml:"cmd"`
}

func GetLocalConf(targetDir string) (*LocalConf, error) {
	yamlFile, err := os.ReadFile(fmt.Sprintf("%s/bib.yml", targetDir))
	if err != nil {
		return nil, fmt.Errorf("cannot read local config file for target directory %s: %v", targetDir, err)
	}
	var c *LocalConf
	if err := yaml.Unmarshal(yamlFile, &c); err != nil {
		return nil, fmt.Errorf("cannot decode YAML contents of local config file for target directory %s: %v", targetDir, err)
	}
	return c, nil
}

func GetFilter(targetDir string, name string) (*Filter, error) {
	c, err := GetLocalConf(targetDir)
	if err != nil {
		return nil, err
	}
	for _, f := range c.Filters {
		if f.Name == name {
			return f, nil
		}
	}
	return nil, fmt.Errorf("target directory %s does not include a filter named '%s'", targetDir, name)
}

func GetTemplates(targetDir string) ([]string, error) {
	files, err := os.ReadDir(fmt.Sprintf("%s/tpl", targetDir))
	if err != nil {
		return nil, fmt.Errorf("cannot read templates for target directory %s: %v", targetDir, err)
	}
	names := []string{}
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}
