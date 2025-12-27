#!/usr/bin/env bash
set -euo pipefail

LOG_DIRS=(
  "./logs"
  "./op-node/logs"
  "./reth-data/logs"
)

echo "Cleaning old log files..."

for dir in "${LOG_DIRS[@]}"; do
  if [ -d "$dir" ]; then
    find "$dir" -type f -name "*.log" -mtime +7 -print -delete
  fi
done

echo "Log cleanup complete."
