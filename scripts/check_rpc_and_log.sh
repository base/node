#!/usr/bin/env bash
set -euo pipefail

# Check L1 and L2 JSON-RPC endpoints and log the results to a file.

L1_RPC="${L1_RPC_URL:-http://localhost:8545}"
L2_RPC="${L2_RPC_URL:-http://localhost:9000}"
LOG_FILE="${LOG_FILE:-./rpc_health.log}"
TIMESTAMP=$(date +"%Y-%m-%d %H:%M:%S")

check_rpc() {
  local name="$1"
  local url="$2"
  if curl -s -H "Content-Type: application/json" \
       -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
       "$url" >/dev/null; then
    echo "$TIMESTAMP [$name] OK" >>"$LOG_FILE"
  else
    echo "$TIMESTAMP [$name] ERROR" >>"$LOG_FILE"
  fi
}

check_rpc "L1" "$L1_RPC"
check_rpc "L2" "$L2_RPC"
