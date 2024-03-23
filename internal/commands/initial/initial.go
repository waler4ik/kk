package initial

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"github.com/waler4ik/kk/internal/tmpl"
	"github.com/waler4ik/kk/internal/walk"
)

type Arguments struct {
	ModulePath  string `long:"modulepath" short:"m" description:"golang module path" default:"rename-or-delete-me"`
	ProjectType string `long:"projecttype" short:"p" description:"golang project type e.g. rest, grpc, graphql" default:"rest"`
}

type Init struct {
	Arguments

	Content embed.FS
}

type TemplateInput struct {
	tmpl.Root
}

func (i *Init) Execute(args []string) error {
	if i.ProjectType == "rest" {
		rootFolder := strings.ToLower(path.Base(i.ModulePath))
		if err := walk.Walk(i.Content, "templates/"+i.ProjectType, rootFolder, &TemplateInput{
			Root: tmpl.Root{
				ModulePath: i.ModulePath,
				Packages:   []tmpl.Package{},
				WSPackages: []tmpl.WSPackage{},
			},
		}); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("currently only REST project type is supported")
	}
	return nil
}
