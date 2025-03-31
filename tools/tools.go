//go:build tools

// This file ensures tool dependencies are kept in sync.  This is the
// recommended way of doing this according to
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// To install the following tools at the version used by this repo run:
// $ make bootstrap
// or
// $ go generate -tags tools tools/tools.go

package tools

//go:generate go install golang.org/x/tools/cmd/goimports
//go:generate go install github.com/mitchellh/gox
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install github.com/twitchtv/twirp/protoc-gen-twirp

import (
	//	_ "golang.org/x/tools/cmd/goimports"
	//
	//	_ "github.com/mitchellh/gox"

	_ "google.golang.org/protobuf/cmd/protoc-gen-go"

	_ "github.com/99designs/gqlgen"

	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
