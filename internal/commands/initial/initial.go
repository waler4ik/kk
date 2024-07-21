package initial

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"github.com/waler4ik/kk/internal/pathchecker"
	"github.com/waler4ik/kk/internal/tmpl"
	"github.com/waler4ik/kk/internal/walk"
)

type Arguments struct {
	RouterType string `long:"router" short:"r" description:"http router e.g. chi, gin" default:"chi"`
	Args       struct {
		ProjectType string `positional-arg-name:"projecttype" description:"golang project type e.g. rest, grpc, graphql" default:"rest"`
		ModulePath  string `positional-arg-name:"modulepath" description:"golang module path" default:"rename-or-delete-me"`
	} `positional-args:"yes" required:"true"`
}

type Init struct {
	Arguments

	Content embed.FS
}

type TemplateInput struct {
	tmpl.Root
}

func (i *Init) Execute(args []string) error {
	if i.Args.ProjectType == "rest" {
		rootFolder := strings.ToLower(path.Base(i.Args.ModulePath))
		templateDir := "templates/" + i.Args.ProjectType
		if err := walk.Walk(i.Content, templateDir, rootFolder, &TemplateInput{
			Root: tmpl.Root{
				ModulePath: i.Args.ModulePath,
				Packages:   []tmpl.Package{},
				WSPackages: []tmpl.WSPackage{},
				RouterType: i.RouterType,
			},
		}, pathchecker.New("", templateDir)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("currently only REST project type is supported")
	}
	return nil
}
