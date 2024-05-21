package main

import "fmt"
import "os"
import "log"
import "gopkg.in/yaml.v3"
import "io/ioutil"

var GLOBAL_CONF = fmt.Sprintf("%s/%s", os.Getenv("XDG_CONFIG_HOME"), "/bib/bib.conf.yml")

type GlobalConf struct {
	Profiles []Profile
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
		} else if subCommand == "list" {
			listProfiles()
		} else if subCommand == "get" {
			getProfile(args)
		}
	}
}

func getProfile(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: bib profiles get <profile>")
		os.Exit(1)
	}
	var name = args[0]
	for _, p := range loadProfiles() {
		if p.Name == name {
			fmt.Println(p.Target)
			os.Exit(0)
		}
	}
	fmt.Fprintf(os.Stderr, "profile '%s' not found\n", name)
	os.Exit(1)
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

func listProfiles() {
	for _, p := range loadProfiles() {
		fmt.Printf("%s\t%s\n", p.Name, p.Target)
	}
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
