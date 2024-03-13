package add

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tangzero/inflector"
	"github.com/waler4ik/kk/internal/walk"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const RESTResource = "resource"

type Config struct {
	Args struct {
		ResourceType string `positional-arg-name:"resourcetype" description:"type of code to add e.g. \"resource\" for a REST resource" required:"yes"`
		Path         string `positional-arg-name:"path" description:"path e.g \"customer\" for the REST resource /customer"`
	} `positional-args:"yes"`
}

type Add struct {
	Config

	Content embed.FS

	ModulePath          string
	ResourceNameSgCaps  string //Singular caps e.g Machine
	ResourceNameSgLower string //Singular lowercase e.g machine
	ResourceNamePlCaps  string //Plural caps e.g Machines
	ResourceNamePlLower string //Plural lowercase e.g machines
}

func (a *Add) Execute(args []string) error {
	if a.Args.ResourceType == RESTResource {
		if a.Args.Path == "" {
			return fmt.Errorf("provide path for the REST resource")
		}
		if !strings.HasPrefix(a.Args.Path, "/") {
			a.Args.Path = "/" + a.Args.Path
		}
		a.prepareResourceNames()

		if err := a.readModulePath(); err != nil {
			return fmt.Errorf("readModulePath: %w", err)
		}

		if err := walk.Walk(a.Content, "templates/"+a.Args.ResourceType, filepath.Clean("internal/endpoints"+a.Args.Path), a); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s resource type not supported", a.Args.ResourceType)
	}

	return nil
}

func (a *Add) readModulePath() error {
	goMod := "go.mod"
	if _, err := os.Stat(goMod); err == nil {
		if goModFile, err := os.ReadFile(goMod); err != nil {
			return fmt.Errorf("ReadFile %s: %s", goMod, err)
		} else {
			a.ModulePath = modfile.ModulePath(goModFile)
		}
	} else {
		return fmt.Errorf("stat: %w", err)
	}
	return nil
}

func (a *Add) prepareResourceNames() {
	basePath := path.Base(a.Args.Path)
	basePathPl := inflector.Pluralize(basePath)
	a.ResourceNamePlCaps = cases.Title(language.English).String(basePathPl)
	a.ResourceNamePlLower = strings.ToLower(basePathPl)
	basePathSg := inflector.Singularize(basePath)
	a.ResourceNameSgCaps = cases.Title(language.English).String(basePathSg)
	a.ResourceNameSgLower = strings.ToLower(basePathSg)
}
