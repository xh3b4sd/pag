# pag

Opinionated protocol buffer api code generator. `pag` stands for "proto api
generator". The tool wraps the necessary business logic of generating [gRPC]
code based on a [protocol buffer] API scheme.



### Prerequisites

`pag` generates `protoc` commands and executes them on the host machine on which
`pag` is executed. In order to generate gRPC code for golang `protoc` and the
golang plugins need to be installed. Depending on the host machine installation
of prerequisites may differ.

```
brew install protobuf
```

```
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```



[gRPC]: https://grpc.io
[protocol buffer]: https://developers.google.com/protocol-buffers
