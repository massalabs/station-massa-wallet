#!/bin/bash -ex

PLUGIN=wallet-plugin

mkdir -p api/server
go generate ./...

# Build the binary
mkdir -p "build/$PLUGIN/"

CGO_LDFLAGS="-framework UniformTypeIdentifiers" CGO_ENABLED="1" go build -tags desktop,production -ldflags "-w -s" -o build/$PLUGIN/wallet-plugin main.go
