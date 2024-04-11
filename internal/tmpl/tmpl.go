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

const (
	Chi = "chi"
	Gin = "gin"
)

type Root struct {
	ModulePath      string
	RouterType      string
	Packages        []Package
	WSPackages      []WSPackage
	ProviderStructs []StructInfo
	SecretManagers  []StructInfo
}
