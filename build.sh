#!/bin/bash
set -ex
# Set the GOPATH environment variable
PLUGIN=wallet-plugin


# Build the binary
mkdir -p "build/$PLUGIN/"
go build -o build/$PLUGIN/thyra-plugin-wallet -ldflags "-X main.port=$1 -X main.path=$2"  cmd/massa-wallet/main.go
