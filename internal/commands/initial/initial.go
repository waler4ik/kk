package initial

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"github.com/waler4ik/kk/internal/walk"
)

type Config struct {
	ModuleName  string `long:"modulename" short:"m" description:"golang module name" default:"rename-or-delete-me"`
	ProjectType string `long:"projecttype" short:"p" description:"golang project type e.g. rest, grpc, graphql" default:"rest"`
}

type Init struct {
	Config

	Content embed.FS
}

func (i *Init) Execute(args []string) error {
	if i.ProjectType == "rest" {
		rootFolder := strings.ToLower(path.Base(i.ModuleName))
		if err := walk.Walk(i.Content, "templates/"+i.ProjectType, rootFolder, i.Config); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("currently only REST project type is supported")
	}
	return nil
}