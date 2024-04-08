package tmpl

type Package struct {
	Name    string
	RelPath string
}

type WSPackage struct {
	Package
	RoutePath string
	NameCaps  string
}

type StructInfo struct {
	Package
	StructName string
}

type Root struct {
	ModulePath      string
	Packages        []Package
	WSPackages      []WSPackage
	ProviderStructs []StructInfo
	SecretManagers  []StructInfo
}
