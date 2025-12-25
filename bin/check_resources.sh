#!/usr/bin/env bash
set -euo pipefail

# Inspect memory and disk usage relevant to a running Base node.
# Useful to determine if your host meets minimum requirements or if cleanup is needed.

echo "== Base node resources =="

# Check available memory (Linux/macOS)
if command -v free >/dev/null 2>&1; then
  free -h
else
  echo "free command not found; skipping memory check"
fi

echo
# Check disk usage under data directories
DATA_DIRS=(
  "${BASE_RETH_DATA_DIR:-./reth-data}"
  "${BASE_GETH_DATA_DIR:-./geth-data}"
)

for dir in "${DATA_DIRS[@]}"; do
  if [ -d "$dir" ]; then
    du -sh "$dir"
  fi
done

echo
df -h .
echo
echo "Consider pruning or moving data directories if disk space is low."
