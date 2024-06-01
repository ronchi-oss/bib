package conf

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type GlobalConf struct {
	Profiles []*Profile
}

type Profile struct {
	Name      string `yaml:"name"`
	TargetDir string `yaml:"target_dir"`
}

func GetGlobalConfPath() (string, error) {
	path := os.Getenv("BIB_GLOBAL_CONFIG")
	if len(path) > 0 {
		if _, err := os.Stat(path); err != nil {
			return "", fmt.Errorf("BIB_GLOBAL_CONFIG is set but does not point to a readable file. Solution: either unset it or point it to a readable file.")
		}
		return path, nil
	}
	if len(os.Getenv("XDG_CONFIG_HOME")) == 0 {
		return "", fmt.Errorf("XDG_CONFIG_HOME is not set so can't use fallback value XDG_CONFIG_HOME/bib/bib.yml as path to global config. Solution: either set BIB_GLOBAL_CONFIG  to a readable file or set XDG_CONFIG_HOME to a readable directory")
	}
	path = fmt.Sprintf("%s/bib/bib.yml", os.Getenv("XDG_CONFIG_HOME"))
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("The fallback global config file %s does not exist. Solution: create it by hand or run `bib create global-config`", path)
	}
	return path, nil
}

func GetGlobalConf() (*GlobalConf, error) {
	path, err := GetGlobalConfPath()
	if err != nil {
		return nil, err
	}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read global config file: %v", err)
	}
	var c *GlobalConf
	if err := yaml.Unmarshal(yamlFile, &c); err != nil {
		return nil, fmt.Errorf("cannot unmarshal data: %v", err)
	}
	return c, nil
}

func GetProfile(name string) (*Profile, error) {
	c, err := GetGlobalConf()
	if err != nil {
		return nil, err
	}
	var p *Profile
	for _, p = range c.Profiles {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("profile '%s' not found", name)
}

func AppendProfile(name string, targetDir string) error {
	c, err := GetGlobalConf()
	if err != nil {
		return err
	}
	for _, p := range c.Profiles {
		if p.Name == name {
			return fmt.Errorf("profile '%s' already exists", p.Name)
		}
	}
	p := &Profile{Name: name, TargetDir: targetDir}
	c.Profiles = append(c.Profiles, p)
	return WriteGlobalConf(c)
}

func WriteGlobalConf(c *GlobalConf) error {
	path, err := GetGlobalConfPath()
	if err != nil {
		return err
	}
	d, err := YAMLEncode(&c)
	if err != nil {
		return fmt.Errorf("cannot marshal global conf: %v", err)
	}
	if err := ioutil.WriteFile(path, d, 0644); err != nil {
		return fmt.Errorf("cannot write global conf: %v", err)
	}
	return nil
}
