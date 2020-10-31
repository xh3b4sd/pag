package typescript

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/pag/pkg/generate"
	"github.com/xh3b4sd/pag/pkg/generate/typescript"
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
		c := typescript.Config{
			FileSystem: afero.NewOsFs(),

			Destination: r.flag.Destination,
			Source:      r.flag.Source,
		}

		g, err = typescript.New(c)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		l, err := g.Commands()
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

			out, err := exec.Command(c.Binary, c.Arguments...).CombinedOutput()
			if err != nil {
				return tracer.Maskf(commandExecutionFailedError, "%s", out)
			}
		}
	}

	{
		l, err := g.Files()
		if err != nil {
			return tracer.Mask(err)
		}

		for _, f := range l {
			// The generated files may define arbitrary file paths on the file
			// system. In order to be super save we simply ensure that the
			// directory in which the generated file is supposed to be written
			// to exists.
			err := os.MkdirAll(filepath.Dir(f.Path), os.ModePerm)
			if err != nil {
				return tracer.Mask(err)
			}

			err = ioutil.WriteFile(f.Path, f.Bytes, 0600)
			if err != nil {
				return tracer.Mask(err)
			}
		}
	}

	return nil
}
