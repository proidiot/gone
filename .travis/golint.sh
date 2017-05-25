#!/bin/sh

set -eu

golint ./... | awk '
/./ {
	if (!found) {
		found = 1
		print "The following issues were reported by golint:"
	}
	print
}
END {
	if (found) {
		exit 1
	} else {
		print "No golint issues."
	}
}'

