package main

import (
	"github.com/ronchi-oss/bib/cmd"
	_ "github.com/ronchi-oss/bib/cmd/cat"
	_ "github.com/ronchi-oss/bib/cmd/create"
	_ "github.com/ronchi-oss/bib/cmd/generate"
	_ "github.com/ronchi-oss/bib/cmd/get"
	_ "github.com/ronchi-oss/bib/cmd/show"
)

func main() {
	cmd.Execute()
}
