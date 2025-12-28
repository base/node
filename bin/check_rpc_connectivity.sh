#!/usr/bin/env bash
set -euo pipefail

# Check connectivity to L1 and OP node JSON-RPC endpoints.
# Usage: BASE_L1_RPC_URL=http://localhost:8545 BASE_OP_RPC_URL=http://localhost:9000 ./bin/check_rpc_connectivity.sh

L1_URL="${BASE_L1_RPC_URL:-http://localhost:8545}"
OP_URL="${BASE_OP_RPC_URL:-http://localhost:9000}"

check_rpc() {
  local url="$1"
  local name="$2"
  echo "Checking $name RPC at $url ..."
  if curl -s -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' "$url" >/dev/null; then
    echo "  $name RPC is reachable."
  else
    echo "  Failed to reach $name RPC."
  fi
}

check_rpc "$L1_URL" "L1"
check_rpc "$OP_URL" "OP node"
