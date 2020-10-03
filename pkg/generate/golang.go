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

func (g *Golang) Generate() ([]Context, error) {
	var dst string
	{
		dst = g.destination
		if strings.HasPrefix(dst, "./") {
			dst = dst[2:]
		}
	}

	dirs := map[string][]string{}
	{
		walkFunc := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return tracer.Mask(err)
			}

			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			if info.IsDir() && info.Name() == ".github" {
				return filepath.SkipDir
			}

			// We do not want to track directories. We are interested in
			// directories containing specific files.
			if info.IsDir() {
				return nil
			}

			// We do not want to track files with the wrong extension. We are
			// interested in protocol buffer files having the .proto extension.
			if filepath.Ext(info.Name()) != ".proto" {
				return nil
			}

			dirs[filepath.Dir(path)] = append(dirs[filepath.Dir(path)], filepath.Join(filepath.Dir(path), info.Name()))

			return nil
		}

		err := afero.Walk(g.fileSystem, g.source, walkFunc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var ctxs []Context
	for d, l := range dirs {
		c := Context{
			Binary:    Binary,
			Arguments: strings.Split(fmt.Sprintf(ArgsFmt, filepath.Join(dst, d), d, strings.Join(l, " ")), " "),
			Directory: filepath.Join(dst, d),
		}

		ctxs = append(ctxs, c)
	}

	return ctxs, nil
}
