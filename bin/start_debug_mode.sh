#!/usr/bin/env bash
set -euo pipefail

# Starts the Base node with verbose logging and debug flags.
# Only recommended for local testing.

export LOG_LEVEL=debug
export OP_NODE_LOG_LEVEL=debug

echo "Starting Base node with debug logs..."
docker compose down
docker compose up --detach

echo "Node is running in debug mode. Use 'docker compose logs -f' to view logs."
