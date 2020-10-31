package typescript

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
	// JsArg is the specific argument string required in order to generate
	// legacy javascript code based on a gRPC api schema. As of time of writing
	// this it is still necessary to generate this code via two separate steps
	// since the upstream ecosystem is migrating towards typescript code
	// generation.
	JsArg = "--js_out=import_style=commonjs,binary:%s/ --proto_path=%s/ %s"
	// TsArg is the specific argument string required in order to generate
	// typescript code based on a gRPC api schema. As of time of writing this it
	// is still necessary to generate this code via two separate steps since the
	// upstream ecosystem is migrating towards typescript code generation.
	TsArg = "--grpc-web_out=import_style=typescript,mode=grpcwebtext:%s/ --proto_path=%s/ %s"
)

type Config struct {
	FileSystem afero.Fs

	Destination string
	Source      string
}

type Typescript struct {
	fileSystem afero.Fs

	destination string
	source      string
}

func New(config Config) (*Typescript, error) {
	if config.FileSystem == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.FileSystem must not be empty", config)
	}

	if config.Destination == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Destination must not be empty", config)
	}
	if config.Source == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Source must not be empty", config)
	}

	t := &Typescript{
		fileSystem: config.FileSystem,

		destination: config.Destination,
		source:      config.Source,
	}

	return t, nil
}

func (t *Typescript) Commands() ([]generate.Command, error) {
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

		err := afero.Walk(t.fileSystem, t.source, walkFunc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var cmds []generate.Command
	for d, l := range dir {
		c := func(f string) generate.Command {
			return generate.Command{
				Binary:    Binary,
				Arguments: strings.Split(fmt.Sprintf(f, filepath.Join(t.destination, d), d, strings.Join(l, " ")), " "),
				Directory: filepath.Join(t.destination, d),
			}
		}

		cmds = append(cmds, c(JsArg))
		cmds = append(cmds, c(TsArg))
	}

	return cmds, nil
}

func (t *Typescript) Files() ([]generate.File, error) {
	return nil, nil
}
