#!/bin/sh

set -eu

if [ "x`golint ./...`" = "x" ]
then
	echo 'Go code passed the linter! Hooray!' >&2
else
	echo 'Go code does not pass lint. Please run: golint ./...' >&2
	exit 1
fi

