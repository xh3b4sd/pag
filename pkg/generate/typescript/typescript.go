package typescript

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

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
	JsArg = "--experimental_allow_proto3_optional --js_out=import_style=commonjs,binary:%s --proto_path=%s %s"
	// TsArg is the specific argument string required in order to generate
	// typescript code based on a gRPC api schema. As of time of writing this it
	// is still necessary to generate this code via two separate steps since the
	// upstream ecosystem is migrating towards typescript code generation.
	TsArg = "--experimental_allow_proto3_optional --grpc-web_out=import_style=typescript,mode=grpcwebtext:%s --proto_path=%s %s"
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
	dirs, err := t.dirs(".proto")
	if err != nil {
		return nil, tracer.Mask(err)
	}

	var cmds []generate.Command
	for _, l := range dirs {
		c := func(f string) generate.Command {
			return generate.Command{
				Binary:    Binary,
				Arguments: strings.Split(fmt.Sprintf(f, t.destination, t.source, strings.Join(l, " ")), " "),
				Directory: t.destination,
			}
		}

		cmds = append(cmds, c(JsArg))
		cmds = append(cmds, c(TsArg))
	}

	return cmds, nil
}

func (t *Typescript) Files() ([]generate.File, error) {
	d, err := t.dirs(".proto")
	if err != nil {
		return nil, tracer.Mask(err)
	}

	var l []generate.File
	{
		p := filepath.Join(t.destination, "index.ts")

		b, err := t.render(p, indexTemplate, t.data(d))
		if err != nil {
			return nil, tracer.Mask(err)
		}

		f := generate.File{
			Path:  p,
			Bytes: b,
		}

		l = append(l, f)
	}

	return l, nil
}

func (t *Typescript) data(dirs map[string][]string) interface{} {
	type Data struct {
		Dir string
	}

	var data []Data

	for d := range dirs {
		data = append(data, Data{Dir: d})
	}

	sort.Slice(data, func(i, j int) bool { return data[i].Dir < data[j].Dir })

	return data
}

func (t *Typescript) dirs(exts ...string) (map[string][]string, error) {
	dirs := map[string][]string{}
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
			for _, e := range exts {
				if filepath.Ext(i.Name()) != e {
					return nil
				}
			}

			dirs[filepath.Dir(p)] = append(dirs[filepath.Dir(p)], filepath.Join(filepath.Dir(p), i.Name()))

			return nil
		}

		err := afero.Walk(t.fileSystem, t.source, walkFunc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return dirs, nil
}

func (t *Typescript) render(path string, tmpl string, data interface{}) ([]byte, error) {
	f := template.FuncMap{
		"ToResource": func(s string) string {
			n := s

			n = filepath.Base(n)
			n = strings.Title(n)

			return n
		},
	}

	s, err := template.New(path).Funcs(f).Parse(tmpl)
	if err != nil {
		return nil, tracer.Mask(err)
	}

	var b bytes.Buffer
	err = s.ExecuteTemplate(&b, path, data)
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return b.Bytes(), nil
}
