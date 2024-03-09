package config

type Config struct {
	ModuleName  string `long:"modulename" description:"golang module name" required:"false" default:"RenameOrDeleteMe"`
	ProjectType string `long:"projecttype" description:"golang project typo e.g. rest, grpc, graphql" required:"false" default:"rest"`
}
