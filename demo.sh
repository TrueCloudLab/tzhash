#!/usr/bin/env bash
#set -x

BLOCK_SIZE=${1:-100M} # 100Mb by default
TMPDIR="${TMPDIR:-$(mktemp -d)}"

OUT="${OUT:-"${TMPDIR}/bighash"}"

echo "Preparing big file at ${OUT}..."
dd if=/dev/urandom of="$OUT" bs="$BLOCK_SIZE" count=1

echo "Make 4 smaller parts from ${OUT}..."
split -dn 4 "${OUT}" "${TMPDIR}/"

echo -n "Big file hash: "
TZALL=$(./bin/tzsum -impl avx2 -name "${OUT}" | awk '{print $1}')
echo "${TZALL}"

for i in $(seq -f "%02g" 0 3)
do
  echo -n "Part ${i} hash:  "
  PART=$(./bin/tzsum -impl avx2 -name "${TMPDIR}/${i}" | awk '{print $1}')
  echo "${PART}" | tee -a "${TMPDIR}/part.hashes"
done

echo -n "Cumulative:    "
TZCUM=$(./bin/homo -concat -file "${TMPDIR}/part.hashes")
echo "${TZCUM}"

if [[ "$TZCUM" == "$TZALL" ]]; then
    echo "Original and cumulative hashes are equal!"
else
    echo "Original and cumulative hashes are NOT equal!"
fi

echo -ne "Cleaning up .. "
rm -rf "${TMPDIR}"
echo "Done!"
