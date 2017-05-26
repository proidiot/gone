#!/bin/sh

set -eu

(
for f in `find . -name '*.go'`
do
	gofmt -l -s ${f}
done
) | awk '
/./ {
	if (!found) {
		found = 1
		print "The following files have gofmt issues:"
	}
	print
}
END {
	if (found) {
		exit 1
	} else {
		print "No gofmt issues."
	}
}'

