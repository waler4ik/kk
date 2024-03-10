package config

type Init struct {
	ModuleName  string `long:"modulename" short:"m" description:"golang module name" required:"false" default:"rename-or-delete-me"`
	ProjectType string `long:"projecttype" short:"p" description:"golang project typo e.g. rest, grpc, graphql" required:"false" default:"rest"`
}
