package add

import "fmt"

type Add struct {
	ArtifactType string `long:"type" short:"t" description:"type of code to add e.g. resource for a REST resource" default:"resource"`
	ArtifactURI  string `long:"uri" short:"u" description:"artifact uri e.g customer for a customer REST resource" default:"default"`
}

func (a *Add) Execute(args []string) error {
	return fmt.Errorf("not implemented")
}
