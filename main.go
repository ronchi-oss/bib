package main

import "fmt"
import "os"
import "log"
import "gopkg.in/yaml.v3"

var BIB_CONF = fmt.Sprintf("%s/%s", os.Getenv("XDG_CONFIG_HOME"), "/bib/bib.conf.yml")

type ConfFile struct {
	Profiles []Profile
}

type Profile struct {
	Name string
	Target string
}

func main() {
	listProfiles()
}

func listProfiles() {
	yamlFile, err := os.ReadFile(BIB_CONF)
	if err != nil {
		log.Fatalf("cannot read global config file: %v", err)
	}
	var conf ConfFile
	if err := yaml.Unmarshal(yamlFile, &conf); err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	for _, p := range conf.Profiles {
		fmt.Printf("%s\t%s\n", p.Name, p.Target)
	}
}
