package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

type GolangConfig struct {
	FileSystem afero.Fs

	Destination string
	Source      string
}

type Golang struct {
	fileSystem afero.Fs

	destination string
	source      string
}

func NewGolang(config GolangConfig) (*Golang, error) {
	if config.FileSystem == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.FileSystem must not be empty", config)
	}

	if config.Destination == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Destination must not be empty", config)
	}
	if config.Source == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Source must not be empty", config)
	}

	g := &Golang{
		fileSystem: config.FileSystem,

		destination: config.Destination,
		source:      config.Source,
	}

	return g, nil
}

func (g *Golang) Generate() ([]Command, error) {
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

	var cmds []Command
	for d, l := range dir {
		c := func(f string) Command {
			return Command{
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
