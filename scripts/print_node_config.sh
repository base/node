#!/usr/bin/env bash
set -euo pipefail

# Small helper script to print the current Base node configuration
# derived from environment variables and defaults.

CLIENT="${CLIENT:-geth}"
HOST_DATA_DIR="${HOST_DATA_DIR:-./${CLIENT}-data}"

echo "Base node configuration:"
echo "  Execution client: ${CLIENT}"
echo "  Host data directory: ${HOST_DATA_DIR}"

if [ -n "${OP_NODE_RPC_URL:-}" ]; then
  echo "  OP_NODE_RPC_URL: ${OP_NODE_RPC_URL}"
else
  echo "  OP_NODE_RPC_URL: not set"
fi

if [ -n "${L1_RPC_URL:-}" ]; then
  echo "  L1_RPC_URL: ${L1_RPC_URL}"
else
  echo "  L1_RPC_URL: not set"
fi
