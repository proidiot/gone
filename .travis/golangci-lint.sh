#!/bin/sh

set -eu

go get -t -u -v github.com/golangci/golangci-lint/cmd/golangci-lint

echo "Running golangci-lint..."
golangci-lint run -v ./...
#golangci-lint run -v --enable-all ./...

echo "No golangci-lint issues found."

