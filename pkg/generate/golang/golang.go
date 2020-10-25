package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/pag/pkg/generate"
)

const (
	Binary = "protoc"
	// MsgArg is the specific argument string required in order to generate go
	// structs based on gRPC messages as of time of writing this. The code
	// generation process is separate due to API changes and migration efforts
	// in the upstream gRPC ecosystem.
	MsgArg = "--go_out=%s/ --proto_path=%s/ %s"
	// SvcArg is the specific argument string required in order to generate go
	// interfaces based on gRPC services as of time of writing this. The code
	// generation process is separate due to API changes and migration efforts
	// in the upstream gRPC ecosystem.
	SvcArg = "--go-grpc_out=%s/ --proto_path=%s/ %s"
)

type Config struct {
	FileSystem afero.Fs

	Destination string
	Source      string
}

type Generate struct {
	fileSystem afero.Fs

	destination string
	source      string
}

func New(config Config) (*Generate, error) {
	if config.FileSystem == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.FileSystem must not be empty", config)
	}

	if config.Destination == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Destination must not be empty", config)
	}
	if config.Source == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Source must not be empty", config)
	}

	g := &Generate{
		fileSystem: config.FileSystem,

		destination: config.Destination,
		source:      config.Source,
	}

	return g, nil
}

func (g *Generate) Generate() ([]generate.Command, error) {
	dir := map[string][]string{}
	{
		walkFunc := func(p string, i os.FileInfo, err error) error {
			if err != nil {
				return tracer.Mask(err)
			}

			if i.IsDir() && i.Name() == ".git" {
				return filepath.SkipDir
			}

			if i.IsDir() && i.Name() == ".github" {
				return filepath.SkipDir
			}

			// We do not want to track directories. We are interested in
			// directories containing specific files.
			if i.IsDir() {
				return nil
			}

			// We do not want to track files with the wrong extension. We are
			// interested in protocol buffer files having the ".proto"
			// extension.
			if filepath.Ext(i.Name()) != ".proto" {
				return nil
			}

			dir[filepath.Dir(p)] = append(dir[filepath.Dir(p)], filepath.Join(filepath.Dir(p), i.Name()))

			return nil
		}

		err := afero.Walk(g.fileSystem, g.source, walkFunc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var cmds []generate.Command
	for d, l := range dir {
		c := func(f string) generate.Command {
			return generate.Command{
				Binary:    Binary,
				Arguments: strings.Split(fmt.Sprintf(f, filepath.Join(g.destination, d), d, strings.Join(l, " ")), " "),
				Directory: filepath.Join(g.destination, d),
			}
		}

		cmds = append(cmds, c(MsgArg))
		cmds = append(cmds, c(SvcArg))
	}

	return cmds, nil
}
