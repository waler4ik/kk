package main

import (
	"embed"
	"log"
	"os"
	"path"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/waler4ik/kk/internal/config"
	"github.com/waler4ik/kk/internal/initial"
)

//go:embed templates
var content embed.FS

func main() {
	config := config.Config{}
	parser := flags.NewParser(&config, flags.Default)
	parser.ShortDescription = "kk - a Golang tool"
	parser.LongDescription = "kk - easily create, setup and extend a Golang projects. No more copying files and endless renaming."
	if _, err := parser.Parse(); err != nil {
		log.Fatalf("Parsing arguments failed: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Provide more arguments. Example: kk init --modulename github.com/yourworkspace/CoolModuleName")
	} else if os.Args[1] == "init" {
		if config.ProjectType == "rest" {
			rootFolder := strings.ToLower(path.Base(config.ModuleName))
			if err := initial.Walk(content, "templates/"+config.ProjectType, rootFolder, config); err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Fatalf("Currently only \"rest\" project type supported")
		}
	} else {
		log.Fatalf("Currently only \"init\" command supported")
	}
}
