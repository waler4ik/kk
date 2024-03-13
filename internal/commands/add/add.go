package add

import (
	"embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
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
const WiringRouterFunctionName = "ConfigureRouter"

type Config struct {
	Args struct {
		ResourceType string `positional-arg-name:"resourcetype" description:"type of code to add e.g. \"resource\" for a REST resource" required:"yes"`
		Path         string `positional-arg-name:"path" description:"path e.g \"customer\" for the REST resource /customer"`
	} `positional-args:"yes"`
}

type Package struct {
	Name    string
	RelPath string
}

type Add struct {
	Config

	Content embed.FS

	ModulePath          string
	PathBase            string
	ResourceNameSgCaps  string //Singular caps e.g Machine
	ResourceNameSgLower string //Singular lowercase e.g machine
	ResourceNamePlCaps  string //Plural caps e.g Machines
	ResourceNamePlLower string //Plural lowercase e.g machines

	Packages []Package
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

		if err := filepath.Walk("internal/endpoints",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				fset := token.NewFileSet()

				if info.IsDir() {
					if pkgMap, err := parser.ParseDir(fset, path,
						func(fi fs.FileInfo) bool {
							return fi.Name() == "controller.go"
						},
						0); err != nil {
						return fmt.Errorf("ParseDir: %w", err)
					} else {
						for _, pkg := range pkgMap {
							for _, file := range pkg.Files {
								for _, decl := range file.Decls {
									if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.String() == WiringRouterFunctionName {
										a.Packages = append(a.Packages, Package{
											Name:    pkg.Name,
											RelPath: path,
										})
									}
								}
							}
						}
					}
				}

				return nil
			}); err != nil {
			return fmt.Errorf("filepath walk: %w", err)
		}
		if err := walk.Walk(a.Content, "templates/rest/internal/config", "internal/config", a); err != nil {
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
	pathBase := path.Base(a.Args.Path)
	a.PathBase = pathBase
	pathBasePl := inflector.Pluralize(pathBase)
	a.ResourceNamePlCaps = cases.Title(language.English).String(pathBasePl)
	a.ResourceNamePlLower = strings.ToLower(pathBasePl)
	pathBaseSg := inflector.Singularize(pathBase)
	a.ResourceNameSgCaps = cases.Title(language.English).String(pathBaseSg)
	a.ResourceNameSgLower = strings.ToLower(pathBaseSg)
}
