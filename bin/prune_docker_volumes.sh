#!/usr/bin/env bash
set -euo pipefail

# Removes all dangling and unused docker volumes to free up space.
# WARNING: this will permanently delete unused volumes.

echo "Pruning unused Docker volumes..."
docker volume prune --force
echo "Done."
