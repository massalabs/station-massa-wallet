#!/bin/bash
set -ex
# Set the GOPATH environment variable
PLUGIN=wallet-plugin


# Build the binary
mkdir -p "build/$PLUGIN/"
CGO_ENABLED="1" go build -o build/$PLUGIN/thyra-plugin-wallet  cmd/massa-wallet/main.go
