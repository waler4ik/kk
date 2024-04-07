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
	"github.com/waler4ik/kk/internal/tmpl"
	"github.com/waler4ik/kk/internal/walk"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const RESTResource = "resource"
const WSResource = "ws"
const EnvSecretManager = "envsecretmanager"
const Postgres = "postgres"
const WiringRouterFunctionName = "ConfigureRouter"
const WiringRouterFunctionContainingFile = "controller.go"
const WebsocketPathConstName = "WsPath"

type Arguments struct {
	Args struct {
		ResourceType string `positional-arg-name:"resourcetype" description:"type of code to add e.g. \"resource\" for a REST resource" required:"yes"`
		Path         string `positional-arg-name:"path" description:"path e.g \"customer\" for the REST resource /customer"`
	} `positional-args:"yes"`
}

type Add struct {
	Arguments

	Content embed.FS
}

type TemplateInput struct {
	tmpl.Root
	PathBase            string
	ResourceNameSgCaps  string //Singular caps e.g Machine
	ResourceNameSgLower string //Singular lowercase e.g machine
	ResourceNamePlCaps  string //Plural caps e.g Machines
	ResourceNamePlLower string //Plural lowercase e.g machines
	RoutePath           string
}

func (a *Add) Execute(args []string) error {

	ti := &TemplateInput{}

	if modPath, err := readModulePath(); err != nil {
		return fmt.Errorf("readModulePath: %w", err)
	} else {
		ti.ModulePath = modPath
	}

	templateDir := "templates/" + a.Args.ResourceType

	if a.Args.ResourceType == RESTResource || a.Args.ResourceType == WSResource {
		if a.Args.Path == "" {
			return fmt.Errorf("path uri is missing")
		}

		if !strings.HasPrefix(a.Args.Path, "/") {
			a.Args.Path = "/" + a.Args.Path
		}

		prepareResourceNames(a.Args.Path, ti)

		ti.RoutePath = a.Args.Path

		if err := walk.Walk(a.Content, templateDir+"/controller", filepath.Clean("internal/endpoints"+a.Args.Path), ti); err != nil {
			return err
		}
		if err := walk.Walk(a.Content, templateDir+"/internal", "internal", ti); err != nil {
			return err
		}
		if err := collectEndpointPackageInfoForWiring("internal/endpoints", ti); err != nil {
			return err
		}

		if err := walk.Walk(a.Content, "templates/rest/internal/config", "internal/config", ti,
			"templates/rest/internal/config/router.go.tmpl",
			"templates/rest/internal/config/websockets.go.tmpl"); err != nil {
			return err
		}

		if err := walk.Walk(a.Content, "templates/rest/internal/api", "internal/api", ti,
			"templates/rest/internal/api/websockets.go.tmpl"); err != nil {
			return err
		}
	} else if a.Args.ResourceType == EnvSecretManager || a.Args.ResourceType == Postgres {
		if err := walk.Walk(a.Content, templateDir, "./", ti); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%s resource type not supported", a.Args.ResourceType)
	}
	return nil
}

func readModulePath() (string, error) {
	goMod := "go.mod"
	if _, err := os.Stat(goMod); err == nil {
		if goModFile, err := os.ReadFile(goMod); err != nil {
			return "", fmt.Errorf("ReadFile %s: %s", goMod, err)
		} else {
			return modfile.ModulePath(goModFile), nil
		}
	} else {
		return "", fmt.Errorf("stat: %w", err)
	}
}

func prepareResourceNames(resourcePath string, ti *TemplateInput) {
	pathBase := path.Base(resourcePath)
	ti.PathBase = pathBase
	pathBasePl := inflector.Pluralize(pathBase)
	ti.ResourceNamePlCaps = cases.Title(language.English).String(pathBasePl)
	ti.ResourceNamePlLower = strings.ToLower(pathBasePl)
	pathBaseSg := inflector.Singularize(pathBase)
	ti.ResourceNameSgCaps = cases.Title(language.English).String(pathBaseSg)
	ti.ResourceNameSgLower = strings.ToLower(pathBaseSg)
}

func collectEndpointPackageInfoForWiring(root string, ti *TemplateInput) error {
	if err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			fset := token.NewFileSet()

			if info.IsDir() {
				if pkgMap, err := parser.ParseDir(fset, path,
					func(fi fs.FileInfo) bool {
						return fi.Name() == WiringRouterFunctionContainingFile
					},
					0); err != nil {
					return fmt.Errorf("ParseDir: %w", err)
				} else {
					for _, pkg := range pkgMap {
						for _, file := range pkg.Files {
							for _, decl := range file.Decls {
								if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.String() == WiringRouterFunctionName {
									ti.Packages = append(ti.Packages, tmpl.Package{
										Name:    pkg.Name,
										RelPath: path,
									})
									return nil
								} else if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok.String() == "const" {
									for _, s := range genDecl.Specs {
										if vs, ok := s.(*ast.ValueSpec); ok {
											for _, name := range vs.Names {
												if name.Name == WebsocketPathConstName {
													for _, value := range vs.Values {
														if bl, ok := value.(*ast.BasicLit); ok {
															ti.WSPackages = append(ti.WSPackages, tmpl.WSPackage{
																Package: tmpl.Package{
																	Name:     pkg.Name,
																	RelPath:  path,
																	NameCaps: cases.Title(language.English).String(pkg.Name),
																},
																RoutePath: bl.Value,
															})
															return nil
														}
													}
												}
											}
										}
									}
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
	return nil
}
