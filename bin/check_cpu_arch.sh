#!/usr/bin/env bash
set -euo pipefail

# This helper script prints the CPU architecture detected at runtime.
# It can be useful when debugging issues related to mismatched container images
# or verifying that a local build uses the correct arch (e.g. amd64 vs arm64).

arch="$(uname -m)"

case "$arch" in
  x86_64)
    echo "Detected architecture: x86_64 (amd64)"
    ;;
  aarch64 | arm64)
    echo "Detected architecture: arm64/aarch64"
    ;;
  *)
    echo "Detected architecture: $arch (unrecognized)"
    ;;
esac

echo "Tip: ensure you pull or build the Base node images matching your host architecture."
