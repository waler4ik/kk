package add

import (
	"embed"
	"errors"
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
const NewFunctionName = "New"
const SecretFunctionName = "Secret"
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

type ResourceTemplateInput struct {
	tmpl.Root
	PathBase            string
	ResourceNameSgCaps  string //Singular caps e.g Machine
	ResourceNameSgLower string //Singular lowercase e.g machine
	ResourceNamePlCaps  string //Plural caps e.g Machines
	ResourceNamePlLower string //Plural lowercase e.g machines
	RoutePath           string
}

type ProviderTemplateInput struct {
	tmpl.Root
}

func (a *Add) Execute(args []string) error {
	modPath, err := readModulePath()
	if err != nil {
		return fmt.Errorf("readModulePath: %w", err)
	}

	templateDir := "templates/" + a.Args.ResourceType

	if a.Args.ResourceType == RESTResource || a.Args.ResourceType == WSResource {
		if a.Args.Path == "" {
			return fmt.Errorf("path uri is missing")
		}

		if !strings.HasPrefix(a.Args.Path, "/") {
			a.Args.Path = "/" + a.Args.Path
		}
		ti := &ResourceTemplateInput{}

		ti.ModulePath = modPath

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
	} else if a.Args.ResourceType == EnvSecretManager {
		pti := &ProviderTemplateInput{}
		pti.ModulePath = modPath

		if _, err := os.Stat("internal/provider"); err == nil {
			if err := collectSecretManagerInfoForWiring("internal/provider", pti); err != nil {
				return err
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat: %w", err)
		}

		if len(pti.SecretManagers) > 0 {
			return fmt.Errorf("you already have a secret manager, try to delete it first and then rerun the add command")
		}

		if err := walk.Walk(a.Content, templateDir, "./", pti); err != nil {
			return err
		}

		if err := collectSecretManagerInfoForWiring("internal/provider", pti); err != nil {
			return err
		}
		if err := collectProviderInfoForWiring("internal/provider", pti); err != nil {
			return err
		}

		if err := walk.Walk(a.Content, "templates/rest/internal/api", "internal/api", pti,
			"templates/rest/internal/api/provider.go.tmpl"); err != nil {
			return err
		}

		if err := walk.Walk(a.Content, "templates/rest/internal/config", "internal/config", pti,
			"templates/rest/internal/config/provider.go.tmpl"); err != nil {
			return err
		}
	} else if a.Args.ResourceType == Postgres {
		pti := &ProviderTemplateInput{}
		pti.ModulePath = modPath

		if _, err := os.Stat("internal/provider"); err == nil {
			if err := collectSecretManagerInfoForWiring("internal/provider", pti); err != nil {
				return err
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat: %w", err)
		}
		smCount := len(pti.SecretManagers)
		if smCount == 0 {
			return fmt.Errorf("add a secret manager first e.g kk add envsecretmanager")
		} else if smCount > 1 {
			return fmt.Errorf("more than one secret manager found under providers, please choose only one")
		}

		if err := walk.Walk(a.Content, templateDir, "./", pti); err != nil {
			return err
		}

		if err := collectProviderInfoForWiring("internal/provider", pti); err != nil {
			return err
		}

		if err := walk.Walk(a.Content, "templates/rest/internal/api", "internal/api", pti,
			"templates/rest/internal/api/provider.go.tmpl"); err != nil {
			return err
		}

		if err := walk.Walk(a.Content, "templates/rest/internal/config", "internal/config", pti,
			"templates/rest/internal/config/provider.go.tmpl"); err != nil {
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

func prepareResourceNames(resourcePath string, ti *ResourceTemplateInput) {
	pathBase := path.Base(resourcePath)
	ti.PathBase = pathBase
	pathBasePl := inflector.Pluralize(pathBase)
	ti.ResourceNamePlCaps = cases.Title(language.English).String(pathBasePl)
	ti.ResourceNamePlLower = strings.ToLower(pathBasePl)
	pathBaseSg := inflector.Singularize(pathBase)
	ti.ResourceNameSgCaps = cases.Title(language.English).String(pathBaseSg)
	ti.ResourceNameSgLower = strings.ToLower(pathBaseSg)
}

func collectEndpointPackageInfoForWiring(root string, ti *ResourceTemplateInput) error {
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
																	Name:    pkg.Name,
																	RelPath: path,
																},
																RoutePath: bl.Value,
																NameCaps:  cases.Title(language.English).String(pkg.Name),
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

func collectProviderInfoForWiring(root string, pti *ProviderTemplateInput) error {
	if err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			fset := token.NewFileSet()

			if info.IsDir() {
				if pkgMap, err := parser.ParseDir(fset, path,
					func(fi fs.FileInfo) bool {
						return strings.HasSuffix(fi.Name(), ".go")
					},
					0); err != nil {
					return fmt.Errorf("ParseDir: %w", err)
				} else {
					for _, pkg := range pkgMap {
						for _, file := range pkg.Files {
							for _, decl := range file.Decls {
								if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.String() == NewFunctionName && funcDecl.Type != nil {
									ft := funcDecl.Type
									if ft.Func != token.NoPos && ft.Params != nil && len(ft.Params.List) == 2 {
										hasLoggerParam := false
										hasSecretManagerParam := false
										if pf, ok := ft.Params.List[0].Type.(*ast.FuncType); ok && pf.Params != nil && len(pf.Params.List) == 2 {
											firstParamType, okFirst := pf.Params.List[0].Type.(*ast.Ident)
											var secondParamType *ast.Ident
											var okSecond bool
											if elli, ok := pf.Params.List[1].Type.(*ast.Ellipsis); ok {
												secondParamType, okSecond = elli.Elt.(*ast.Ident)
											}

											if okFirst && okSecond && firstParamType.Name == "string" && secondParamType.Name == "any" {
												hasLoggerParam = true
											}
										} else {
											continue
										}

										if se, ok := ft.Params.List[1].Type.(*ast.SelectorExpr); ok {
											if se.Sel != nil && se.Sel.Name == "SecretManager" {
												if id, ok := se.X.(*ast.Ident); ok && id.Name == "secretmanager" {
													hasSecretManagerParam = true
												}
											}
										} else {
											continue
										}

										if hasLoggerParam && hasSecretManagerParam {
											if ft.Results != nil && len(ft.Results.List) == 2 {
												var firstReturnTypeName, secondReturnTypeName string
												if sx, ok := ft.Results.List[0].Type.(*ast.StarExpr); ok {
													if rt, ok := sx.X.(*ast.Ident); ok {
														firstReturnTypeName = rt.Name
													}
												}
												if sx, ok := ft.Results.List[1].Type.(*ast.Ident); ok {
													secondReturnTypeName = sx.Name
												}
												if firstReturnTypeName != "" && secondReturnTypeName == "error" {
													pti.ProviderStructs = append(pti.ProviderStructs, tmpl.StructInfo{
														Package: tmpl.Package{
															Name:    pkg.Name,
															RelPath: path,
														},
														StructName: firstReturnTypeName,
													})
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

func collectSecretManagerInfoForWiring(root string, pti *ProviderTemplateInput) error {
	if err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			fset := token.NewFileSet()

			if info.IsDir() {
				if pkgMap, err := parser.ParseDir(fset, path,
					func(fi fs.FileInfo) bool {
						return strings.HasSuffix(fi.Name(), ".go")
					},
					0); err != nil {
					return fmt.Errorf("ParseDir: %w", err)
				} else {
					for _, pkg := range pkgMap {
						var structInfo tmpl.StructInfo
						foundSecretMethod := false
						for _, file := range pkg.Files {
							for _, decl := range file.Decls {
								if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Type != nil {
									ft := funcDecl.Type
									if funcDecl.Name.String() == NewFunctionName {
										if ft.Func != token.NoPos && ft.Params != nil || len(ft.Params.List) == 1 {
											if pf, ok := ft.Params.List[0].Type.(*ast.FuncType); ok && pf.Params != nil && len(pf.Params.List) == 2 {
												firstParamType, okFirst := pf.Params.List[0].Type.(*ast.Ident)
												var secondParamType *ast.Ident
												var okSecond bool
												if elli, ok := pf.Params.List[1].Type.(*ast.Ellipsis); ok {
													secondParamType, okSecond = elli.Elt.(*ast.Ident)
												}

												if okFirst && okSecond && firstParamType.Name == "string" && secondParamType.Name == "any" {
													if ft.Results != nil && len(ft.Results.List) > 0 {
														if sx, ok := ft.Results.List[0].Type.(*ast.StarExpr); ok {
															if rt, ok := sx.X.(*ast.Ident); ok {
																structInfo = tmpl.StructInfo{
																	Package: tmpl.Package{
																		Name:    pkg.Name,
																		RelPath: path,
																	},
																	StructName: rt.Name,
																}
															}
														}

													}
												}
											} else {
												continue
											}
										}
									} else if funcDecl.Name.String() == SecretFunctionName && funcDecl.Recv != nil {
										for _, receiver := range funcDecl.Recv.List {
											if sx, ok := receiver.Type.(*ast.StarExpr); ok {
												if rt, ok := sx.X.(*ast.Ident); ok && rt.Name == structInfo.StructName {
													foundSecretMethod = true
												}
											}
										}
									}
								}
							}
						}
						if structInfo.Name != "" && foundSecretMethod {
							pti.SecretManagers = append(pti.SecretManagers, structInfo)
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
