package generate

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

var update = flag.Bool("update", false, "update .golden files")

// TODO
// Test_Generate_Golang tests the masking behaviour based on our Error type's JSON
// output. The tests use golden file references. In case the golden files change
// something is broken. In case intentional changes get introduced the golden
// files have to be updated. In case the golden files have to be adjusted,
// simply provide the -update flag when running the tests.
//
//     go test . -run Test_Generate_Golang -update
//
func Test_Generate_Golang(t *testing.T) {
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
				c := GolangConfig{
					FileSystem: tc.fs,

					Destination: tc.dst,
					Source:      tc.src,
				}

				g, err = NewGolang(c)
				if err != nil {
					t.Fatal(err)
				}
			}

			l, err := g.Generate()
			if err != nil {
				t.Fatal(err)
			}

			var actual string
			for _, c := range l {
				actual += c.String() + "\n"
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

			if !bytes.Equal([]byte(actual), expected) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(expected), string(actual)))
			}
		})
	}
}

func fileName(i int) string {
	return "case-" + strconv.Itoa(i) + ".golden"
}
