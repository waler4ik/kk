package main

import (
	"embed"

	flags "github.com/jessevdk/go-flags"
	"github.com/waler4ik/kk/internal/commands/add"
	"github.com/waler4ik/kk/internal/commands/initial"
)

//go:embed templates
var content embed.FS

func main() {
	parser := flags.NewNamedParser("kk", flags.Default)
	parser.ShortDescription = "kk - a Golang tool"
	parser.LongDescription = "kk - easily create, setup and extend a Golang projects. No more copying files and endless renaming."

	initCMD := initial.Init{Content: content}
	if _, err := parser.AddCommand("init", "Creates a Golang project from internal templates", "Creates a Golang project for REST, GRPC, GraphQL in a separate folder", &initCMD); err != nil {
		panic(err)
	}

	addCMD := add.Add{Content: content}
	if _, err := parser.AddCommand("add", "Adds a resource to current project", "Adds and properly wires e.g REST resource", &addCMD); err != nil {
		panic(err)
	}

	if _, err := parser.Parse(); err != nil {
		return
	}
}
