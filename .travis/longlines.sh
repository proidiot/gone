#!/bin/sh

set -eu

(
for f in `find . -name '*.go'`
do
	cat ${f} | awk -v f="${f}" '
	{
		l = o = 0

		while (n = index(substr($0, o + 1), "\t")) {
			o += n
			l += n + (8 - ((l + n) % 8))
		}

		l += length(substr($0, o + 1))

		if (l > 80) {
			print f
			exit
		}
	}'
done
) | awk '
/./ {
	if (!found) {
		found = 1
		print "The following files have lines too long:"
	}
	print
}
END {
	if (found) {
		exit 1
	} else {
		print "No files were too long."
	}
}'

