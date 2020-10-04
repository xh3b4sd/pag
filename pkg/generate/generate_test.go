package generate

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_Generate_Generate tests the protoc command generation. The protoc binary
// is used to generate language specific code based on a gRPC apischema. The
// generated protoc commands are executed in order to generate the actual
// language specific code. The tests here ensure that the command execution with
// its flags and positional arguments works as expected.
//
//     go test ./... -run Test_Generate_Generate -update
//
func Test_Generate_Generate(t *testing.T) {
	testCases := []struct {
		fs  afero.Fs
		dst string
		src string
	}{
		// Case 0 ensures that a single proto file in a single directory is
		// scanned accordingly.
		{
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()

				fs.MkdirAll("pbf", 0755)
				afero.WriteFile(fs, "pbf/foo.proto", nil, 0644)

				return fs
			}(),
			dst: "./pkg/",
			src: ".",
		},
		// Case 1 ensures that a single proto file in a single directory is
		// scanned accordingly. Note that there are additional empty
		// directories.
		{
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()

				fs.MkdirAll("pbf", 0755)
				afero.WriteFile(fs, "pbf/foo.proto", nil, 0644)

				fs.MkdirAll("pbf/post", 0755)
				fs.MkdirAll("pbf/user", 0755)

				return fs
			}(),
			dst: "some/other/dir",
			src: ".",
		},
		// Case 2 ensures that multiple proto files in multiple directories are
		// scanned accordingly.
		{
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()

				fs.MkdirAll("pbf/post", 0755)
				afero.WriteFile(fs, "pbf/post/api.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/post/create.proto", nil, 0644)

				fs.MkdirAll("pbf/user", 0755)
				afero.WriteFile(fs, "pbf/user/foo.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/user/bar.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/user/baz.proto", nil, 0644)

				return fs
			}(),
			dst: "./pkg/",
			src: ".",
		},
		// Case 3 ensures that multiple proto files in multiple directories are
		// scanned accordingly. Note that the given destination path is slightly
		// different.
		{
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()

				fs.MkdirAll("pbf/post", 0755)
				afero.WriteFile(fs, "pbf/post/api.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/post/create.proto", nil, 0644)

				fs.MkdirAll("pbf/user", 0755)
				afero.WriteFile(fs, "pbf/user/foo.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/user/bar.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/user/baz.proto", nil, 0644)

				return fs
			}(),
			dst: "pkg",
			src: ".",
		},
		// Case 4 ensures that multiple proto files in multiple directories are
		// scanned accordingly. Note that the given destination path is
		// absolute.
		{
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()

				fs.MkdirAll("pbf/post", 0755)
				afero.WriteFile(fs, "pbf/post/api.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/post/create.proto", nil, 0644)

				fs.MkdirAll("pbf/user", 0755)
				afero.WriteFile(fs, "pbf/user/foo.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/user/bar.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/user/baz.proto", nil, 0644)

				fs.MkdirAll("pbf/more/deeply/nested", 0755)
				afero.WriteFile(fs, "pbf/more/deeply/nested/foo.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/more/deeply/nested/bar.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/more/deeply/nested/baz.proto", nil, 0644)

				return fs
			}(),
			dst: "/home/runner/tmp/pkg/",
			src: ".",
		},
		// Case 5 ensures that only proto files in the source directory are
		// scanned accordingly.
		{
			fs: func() afero.Fs {
				fs := afero.NewMemMapFs()

				fs.MkdirAll("pbf/post", 0755)
				afero.WriteFile(fs, "pbf/post/api.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/post/create.proto", nil, 0644)

				fs.MkdirAll("pbf/user", 0755)
				afero.WriteFile(fs, "pbf/user/foo.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/user/bar.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/user/baz.proto", nil, 0644)

				fs.MkdirAll("pbf/more/deeply/nested", 0755)
				afero.WriteFile(fs, "pbf/more/deeply/nested/foo.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/more/deeply/nested/bar.proto", nil, 0644)
				afero.WriteFile(fs, "pbf/more/deeply/nested/baz.proto", nil, 0644)

				return fs
			}(),
			dst: "./pkg/",
			src: "./pbf/user/",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var g Interface
			{
				c := Config{
					FileSystem: tc.fs,

					Destination: tc.dst,
					Source:      tc.src,
				}

				g, err = New(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			l, err := g.Generate()
			if err != nil {
				t.Fatal(err)
			}

			var actual string
			{
				var s []string
				for _, c := range l {
					s = append(s, c.String())
				}

				sort.Strings(s)

				actual = strings.Join(s, "\n")
			}

			p := filepath.Join("testdata/generate", fileName(i))
			if *update {
				err := ioutil.WriteFile(p, []byte(actual), 0644) // nolint:gosec
				if err != nil {
					t.Fatal(err)
				}
			}

			expected, err := ioutil.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expected, []byte(actual)) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(actual), string(expected)))
			}
		})
	}
}

func fileName(i int) string {
	return "case-" + strconv.Itoa(i) + ".golden"
}
