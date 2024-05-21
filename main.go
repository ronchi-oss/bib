package main

import "fmt"
import "gopkg.in/yaml.v3"
import "io/ioutil"
import "log"
import "os"
import "strings"

var GLOBAL_CONF = fmt.Sprintf("%s/%s", os.Getenv("XDG_CONFIG_HOME"), "/bib/bib.go.yml")

type GlobalConf struct {
	Profiles []Profile
}

type NoteFilterSpec struct {
	Name    string
	Pattern string
}

type Profile struct {
	Name   string
	Target string
}

func main() {
	var args = os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: bib <command>")
		os.Exit(1)
	}
	command := args[0]
	args = args[1:]
	if command == "profiles" {
		if len(args) == 0 {
			fmt.Println("Usage: bib profiles <command>")
			os.Exit(1)
		}
		subCommand := args[0]
		args = args[1:]
		if subCommand == "add" {
			addProfile(args)
		} else if subCommand == "delete" {
			deleteProfile(args)
		} else if subCommand == "list" {
			listProfiles()
		} else if subCommand == "get" {
			getProfile(args)
		}
	} else if command == "profile-filters" {
		getProfileFilters(args)
	}
}

func getProfile(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: bib profiles get <profile>")
		os.Exit(1)
	}
	p, err := loadProfile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "profile '%s' not found\n", args[0])
		os.Exit(1)
	}
	fmt.Println(p.Target)
}

func getProfileFilters(args []string) {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: bib profile-filters <profile>\n")
		os.Exit(1)
	}
	p, err := loadProfile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "profile '%s' not found\n", args[0])
		os.Exit(1)
	}
	yamlFile, err := os.ReadFile(getTargetFiltersConfPath(p.Target))
	if err != nil {
		log.Fatalf("cannot read filters file: %v", err)
	}
	var nfSpecs []NoteFilterSpec
	if err := yaml.Unmarshal(yamlFile, &nfSpecs); err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	fmt.Println("all")
	for _, s := range nfSpecs {
		fmt.Println(s.Name)
	}
}

func addProfile(args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: bib profiles add <name> <target>\n")
		os.Exit(1)
	}
	c := loadGlobalConf()
	for _, p := range c.Profiles {
		if p.Name == args[0] {
			fmt.Fprintf(os.Stderr, "profile '%s' already exists\n", p.Name)
			os.Exit(1)
		}
	}
	p := Profile{Name: args[0], Target: args[1]}
	c.Profiles = append(c.Profiles, p)
	writeGlobalConf(c)
}

func deleteProfile(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: bib profiles delete <name>\n")
		os.Exit(1)
	}
	c := loadGlobalConf()
	found := false
	for i, p := range c.Profiles {
		if p.Name == args[0] {
			found = true
			c.Profiles = append(c.Profiles[:i], c.Profiles[i+1:]...)
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "profile '%s' not found\n", args[0])
		os.Exit(1)
	}
	writeGlobalConf(c)
}

func listProfiles() {
	for _, p := range loadProfiles() {
		fmt.Printf("%s\t%s\n", p.Name, p.Target)
	}
}

func loadProfile(name string) (Profile, error) {
	var p Profile
	for _, p = range loadProfiles() {
		if p.Name == name {
			return p, nil
		}
	}
	return p, fmt.Errorf("profile '%s' not found\n", name)
}

func loadProfiles() []Profile {
	return loadGlobalConf().Profiles
}

func loadGlobalConf() GlobalConf {
	yamlFile, err := os.ReadFile(GLOBAL_CONF)
	if err != nil {
		log.Fatalf("cannot read global config file: %v", err)
	}
	var c GlobalConf
	if err := yaml.Unmarshal(yamlFile, &c); err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	return c
}

func writeGlobalConf(c GlobalConf) {
	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("cannot marshal global conf: %v", err)
	}
	if err := ioutil.WriteFile(GLOBAL_CONF, d, 0644); err != nil {
		log.Fatalf("cannot write global conf: %v", err)
	}
}

func getTargetFiltersConfPath(targetDir string) string {
	if strings.HasPrefix(targetDir, "~/") {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		targetDir = userHomeDir + "/" + targetDir[2:]
	}
	return targetDir + "/bib-filters.yml"
}
