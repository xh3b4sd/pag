package golang

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type flag struct {
	Destination string
	Source      string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Destination, "destination", "d", "./pkg/", "Directory to put the generated golang code into.")
	cmd.Flags().StringVarP(&f.Source, "source", "s", ".", "Directory to look for the gRPC api schema definitions.")
}

func (f *flag) Validate() error {
	if f.Destination == "" {
		return tracer.Maskf(invalidFlagError, "-d/--destination must not be empty")
	}
	if f.Source == "" {
		return tracer.Maskf(invalidFlagError, "-s/--source must not be empty")
	}

	return nil
}
