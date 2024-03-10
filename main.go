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
	config := config.Init{}
	parser := flags.NewNamedParser("kk", flags.Default)
	parser.ShortDescription = "kk - a Golang tool"
	parser.LongDescription = "kk - easily create, setup and extend a Golang projects. No more copying files and endless renaming."
	parser.AddCommand("init", "Creates a Golang project from internal templates", "Creates a Golang project for rest, grpc, graphql in separate folder", &config)
	if _, err := parser.Parse(); err != nil {
		return
	}

	if len(os.Args) < 2 {
		parser.WriteHelp(os.Stdout)
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
