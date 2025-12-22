#!/usr/bin/env bash
set -euo pipefail

# Simple JSON-RPC healthcheck for Base op-geth.
#
# Usage:
#   BASE_GETH_RPC_URL=http://localhost:8545 ./geth/check_geth_rpc.sh
#
# The script calls `eth_blockNumber` and considers the node healthy
# if it returns a non-empty result.

RPC_URL="${BASE_GETH_RPC_URL:-http://localhost:8545}"
TIMEOUT="${BASE_GETH_RPC_TIMEOUT:-3}"

log_err() {
  echo "[geth-health] $*" >&2
}

if ! command -v curl >/dev/null 2>&1; then
  log_err "curl not found in PATH"
  exit 1
fi

payload='{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}'

response="$(curl -sS \
  --max-time "${TIMEOUT}" \
  -H "Content-Type: application/json" \
  -d "${payload}" \
  "${RPC_URL}" 2>/dev/null || true)"

if [ -z "${response}" ]; then
  log_err "empty response from ${RPC_URL}"
  exit 1
fi

if ! echo "${response}" | grep -q '"result"'; then
  log_err "eth_blockNumber did not return a result: ${response}"
  exit 1
fi

echo "[geth-health] ok: JSON-RPC is responding at ${RPC_URL}"
exit 0
