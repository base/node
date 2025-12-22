#!/usr/bin/env bash
set -euo pipefail

# Small helper to inspect disk usage of a Base Reth data directory.
#
# This can be used before upgrades or when planning storage for a new node.
#
# Usage:
#   BASE_RETH_DATADIR=/data/reth ./reth/check_reth_datadir.sh

DATADIR="${BASE_RETH_DATADIR:-/data/reth}"

if [ ! -d "${DATADIR}" ]; then
  echo "[reth] Data directory does not exist:"
  echo "  ${DATADIR}"
  echo
  echo "Set BASE_RETH_DATADIR to the correct path before running this script."
  exit 1
fi

echo "== Base Reth datadir disk usage =="
echo "Path: ${DATADIR}"
echo

if command -v du >/dev/null 2>&1; then
  du -sh "${DATADIR}"
  echo
  echo "Top-level subdirectories:"
  du -sh "${DATADIR}"/* 2>/dev/null | sort -h
else
  echo "The 'du' command is not available on this system."
fi

echo
echo "Tip: keep this directory on a fast local NVMe SSD for best performance."
