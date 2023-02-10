#!/bin/bash

BLOCK_SIZE=${1:-1G} # gigabyte by default
OUT="${OUT:-$(mktemp /tmp/random-file.XXXXXX)}"

dd if=/dev/urandom of="$OUT" bs="$BLOCK_SIZE" count=1

for impl in avx avx2 generic; do
	echo $impl implementation:
	time ./bin/tzsum -name "$OUT" -impl $impl
	echo
done
