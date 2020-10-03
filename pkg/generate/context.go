package generate

import (
	"fmt"
	"strings"
)

const (
	Binary  = "protoc"
	ArgsFmt = "--go-grpc_out=./%s/ --proto_path=./%s/ %s"
)

type Context struct {
	// Arguments are the command line arguments appended to the binary used to
	// generate language specific code for a gRPC api schema.
	Arguments []string
	// Binary is the binary used to generate language specific code for a gRPC
	// api schema. Typically it should simply be "protoc".
	Binary string
	// Directory is the file system location in which language specific code
	// will be generated into. This particular information is exposed because
	// the directory structure has to be ensured so that the gRPC tooling can
	// work properly since it is extremely picky with folders not existing.
	Directory string
}

// String joins the binary and its arguments resulting in one concatenated
// string representation of the complete executable command. Note the example
// below and how the ".proto" files are provided as separate positional
// arguments.
//
//     protoc \
//
//         --go-grpc_out=./pkg/pbf/user/ \
//         --proto_path=./pbf/user/ \
//
//         pbf/user/api.proto \
//         pbf/user/create.proto \
//         pbf/user/delete.proto \
//         pbf/user/search.proto \
//         pbf/user/update.proto
//
func (c Context) String() string {
	return fmt.Sprintf("%s %s", c.Binary, strings.Join(c.Arguments, " "))
}
