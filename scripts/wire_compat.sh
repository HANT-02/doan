#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <package_dir> <wire_gen_file>" >&2
  exit 1
fi

package_dir="$1"
wire_gen_file="$2"

if command -v wire >/dev/null 2>&1; then
  if (
    cd "$package_dir"
    wire gen .
  ); then
    exit 0
  fi

  echo "wire gen failed for ${package_dir}; using checked-in ${wire_gen_file}." >&2
fi

if [[ -f "$wire_gen_file" ]]; then
  echo "Reusing existing ${wire_gen_file}. Installed wire is likely incompatible with the current Go toolchain." >&2
  exit 0
fi

echo "wire generation failed and ${wire_gen_file} does not exist." >&2
exit 1
