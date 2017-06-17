#!/bin/sh

set -eu

go get -t -u -v github.com/alecthomas/gometalinter

gometalinter -i -u

go get -t -u -v ./...

gometalinter \
	--enable-all \
	--deadline=600s \
	--enable-gc \
	--cyclo-over=15 \
	./...

