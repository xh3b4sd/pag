package golang

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/pag/pkg/generate"
)

type runner struct {
	flag   *flag
	logger logger.Interface
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return tracer.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	var g generate.Interface
	{
		c := generate.GolangConfig{
			FileSystem: afero.NewOsFs(),

			Destination: r.flag.Destination,
			Source:      r.flag.Source,
		}

		g, err = generate.NewGolang(c)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		l, err := g.Generate()
		if err != nil {
			return tracer.Mask(err)
		}

		for _, c := range l {
			// The gRPC tooling is not particularly prudent with file path and
			// file system management. We need to ensure the configured
			// directory structure in advance so that the gRPC tooling can
			// generate the language specific code into that.
			err := os.MkdirAll(c.Directory, os.ModePerm)
			if err != nil {
				return tracer.Mask(err)
			}

			fmt.Printf("%#v\n", c.String())

			out, err := exec.Command(c.Binary, c.Arguments...).CombinedOutput()
			if err != nil {
				fmt.Printf("%s\n", out)
				return tracer.Mask(err)
			}
		}
	}

	return nil
}
