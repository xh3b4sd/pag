package typescript

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

	"github.com/xh3b4sd/pag/pkg/generate"
)

var update = flag.Bool("update", false, "update .golden files")

// Test_Typescript_Commands tests the protoc command generation. The protoc
// binary is used to generate language specific code based on a gRPC apischema.
// The generated protoc commands are executed in order to generate the actual
// language specific code. The tests here ensure that the command execution with
// its flags and positional arguments works as expected.
//
//     go test ./... -run Test_Typescript_Commands -update
//
func Test_Typescript_Commands(t *testing.T) {
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

				mustCreateDir(fs, "pbf")
				mustCreateFile(fs, "pbf/foo.proto")

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

				mustCreateDir(fs, "pbf")
				mustCreateFile(fs, "pbf/foo.proto")

				mustCreateDir(fs, "pbf/post")
				mustCreateDir(fs, "pbf/user")

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

				mustCreateDir(fs, "pbf/post")
				mustCreateFile(fs, "pbf/post/api.proto")
				mustCreateFile(fs, "pbf/post/create.proto")

				mustCreateDir(fs, "pbf/user")
				mustCreateFile(fs, "pbf/user/foo.proto")
				mustCreateFile(fs, "pbf/user/bar.proto")
				mustCreateFile(fs, "pbf/user/baz.proto")

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

				mustCreateDir(fs, "pbf/post")
				mustCreateFile(fs, "pbf/post/api.proto")
				mustCreateFile(fs, "pbf/post/create.proto")

				mustCreateDir(fs, "pbf/user")
				mustCreateFile(fs, "pbf/user/foo.proto")
				mustCreateFile(fs, "pbf/user/bar.proto")
				mustCreateFile(fs, "pbf/user/baz.proto")

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

				mustCreateDir(fs, "pbf/post")
				mustCreateFile(fs, "pbf/post/api.proto")
				mustCreateFile(fs, "pbf/post/create.proto")

				mustCreateDir(fs, "pbf/user")
				mustCreateFile(fs, "pbf/user/foo.proto")
				mustCreateFile(fs, "pbf/user/bar.proto")
				mustCreateFile(fs, "pbf/user/baz.proto")

				mustCreateDir(fs, "pbf/more/deeply/nested")
				mustCreateFile(fs, "pbf/more/deeply/nested/foo.proto")
				mustCreateFile(fs, "pbf/more/deeply/nested/bar.proto")
				mustCreateFile(fs, "pbf/more/deeply/nested/baz.proto")

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

				mustCreateDir(fs, "pbf/post")
				mustCreateFile(fs, "pbf/post/api.proto")
				mustCreateFile(fs, "pbf/post/create.proto")

				mustCreateDir(fs, "pbf/user")
				mustCreateFile(fs, "pbf/user/foo.proto")
				mustCreateFile(fs, "pbf/user/bar.proto")
				mustCreateFile(fs, "pbf/user/baz.proto")

				mustCreateDir(fs, "pbf/more/deeply/nested")
				mustCreateFile(fs, "pbf/more/deeply/nested/foo.proto")
				mustCreateFile(fs, "pbf/more/deeply/nested/bar.proto")
				mustCreateFile(fs, "pbf/more/deeply/nested/baz.proto")

				return fs
			}(),
			dst: "./pkg/",
			src: "./pbf/user/",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var g generate.Interface
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

			l, err := g.Commands()
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

				actual = strings.Join(s, "\n") + "\n"
			}

			p := filepath.Join("testdata/commands", fileName(i))
			if *update {
				err := ioutil.WriteFile(p, []byte(actual), 0600)
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

// Test_Typescript_Files tests the protoc command generation. The protoc
// binary is used to generate language specific code based on a gRPC apischema.
// The generated protoc commands are executed in order to generate the actual
// language specific code. The tests here ensure that the command execution with
// its flags and positional arguments works as expected.
//
//     go test ./... -run Test_Typescript_Files -update
//
func Test_Typescript_Files(t *testing.T) {
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

				mustCreateDir(fs, "pbf")
				mustCreateFile(fs, "pbf/foo.proto")

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

				mustCreateDir(fs, "pbf")
				mustCreateFile(fs, "pbf/foo.proto")

				mustCreateDir(fs, "pbf/post")
				mustCreateDir(fs, "pbf/user")

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

				mustCreateDir(fs, "pbf/post")
				mustCreateFile(fs, "pbf/post/api.proto")
				mustCreateFile(fs, "pbf/post/create.proto")

				mustCreateDir(fs, "pbf/user")
				mustCreateFile(fs, "pbf/user/foo.proto")
				mustCreateFile(fs, "pbf/user/bar.proto")
				mustCreateFile(fs, "pbf/user/baz.proto")

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

				mustCreateDir(fs, "pbf/post")
				mustCreateFile(fs, "pbf/post/api.proto")
				mustCreateFile(fs, "pbf/post/create.proto")

				mustCreateDir(fs, "pbf/user")
				mustCreateFile(fs, "pbf/user/foo.proto")
				mustCreateFile(fs, "pbf/user/bar.proto")
				mustCreateFile(fs, "pbf/user/baz.proto")

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

				mustCreateDir(fs, "pbf/post")
				mustCreateFile(fs, "pbf/post/api.proto")
				mustCreateFile(fs, "pbf/post/create.proto")

				mustCreateDir(fs, "pbf/user")
				mustCreateFile(fs, "pbf/user/foo.proto")
				mustCreateFile(fs, "pbf/user/bar.proto")
				mustCreateFile(fs, "pbf/user/baz.proto")

				mustCreateDir(fs, "pbf/more/deeply/nested")
				mustCreateFile(fs, "pbf/more/deeply/nested/foo.proto")
				mustCreateFile(fs, "pbf/more/deeply/nested/bar.proto")
				mustCreateFile(fs, "pbf/more/deeply/nested/baz.proto")

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

				mustCreateDir(fs, "pbf/post")
				mustCreateFile(fs, "pbf/post/api.proto")
				mustCreateFile(fs, "pbf/post/create.proto")

				mustCreateDir(fs, "pbf/user")
				mustCreateFile(fs, "pbf/user/foo.proto")
				mustCreateFile(fs, "pbf/user/bar.proto")
				mustCreateFile(fs, "pbf/user/baz.proto")

				mustCreateDir(fs, "pbf/more/deeply/nested")
				mustCreateFile(fs, "pbf/more/deeply/nested/foo.proto")
				mustCreateFile(fs, "pbf/more/deeply/nested/bar.proto")
				mustCreateFile(fs, "pbf/more/deeply/nested/baz.proto")

				return fs
			}(),
			dst: "./pkg/",
			src: "./pbf/user/",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			var g generate.Interface
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

			l, err := g.Files()
			if err != nil {
				t.Fatal(err)
			}

			var actual string
			{
				var s []string
				for _, f := range l {
					s = append(s, f.Path)
					s = append(s, string(f.Bytes))
				}

				sort.Strings(s)

				actual = strings.Join(s, "\n") + "\n"
			}

			p := filepath.Join("testdata/files", fileName(i))
			if *update {
				err := ioutil.WriteFile(p, []byte(actual), 0600)
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

func mustCreateDir(fs afero.Fs, p string) {
	err := fs.MkdirAll(p, 0755)
	if err != nil {
		panic(err)
	}
}

func mustCreateFile(fs afero.Fs, p string) {
	err := afero.WriteFile(fs, p, nil, 0644)
	if err != nil {
		panic(err)
	}
}
