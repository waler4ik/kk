package tmpl

type Package struct {
	Name     string
	RelPath  string
	NameCaps string
}

type WSPackage struct {
	Package
	RoutePath string
}

type Root struct {
	ModulePath string
	Packages   []Package
	WSPackages []WSPackage
}
