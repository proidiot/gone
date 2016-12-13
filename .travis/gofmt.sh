#!/bin/sh

cd `dirname $0`

if [ "x`go fmt ./...`" = "x" ]
then
	echo 'Go code is formatted properly! Hooray!' >&2
else
	echo 'Go code is not formatted properly. Please run: go fmt ./...' >&2
	exit 1
fi

