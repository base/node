#!/usr/bin/env bash
set -euo pipefail

# Simple healthcheck for a Nethermind JSON-RPC endpoint.
# Can be used from Docker HEALTHCHECK or external monitoring.
#
# Environment (optional):
#   BASE_NETHERMIND_RPC_URL    - RPC URL to probe (default: http://localhost:8545)
#   BASE_NETHERMIND_RPC_TIMEOUT - curl timeout in seconds (default: 3)

RPC_URL="${BASE_NETHERMIND_RPC_URL:-http://localhost:8545}"
TIMEOUT="${BASE_NETHERMIND_RPC_TIMEOUT:-3}"

log_err() {
  echo "[nethermind-health] $*" >&2
}

if ! command -v curl >/dev/null 2>&1; then
  log_err "curl not found in PATH"
  exit 1
fi

payload='{"jsonrpc":"2.0","id":1,"method":"eth_syncing","params":[]}'

response="$(curl -sS \
  --max-time "${TIMEOUT}" \
  -H "Content-Type: application/json" \
  -d "${payload}" \
  "${RPC_URL}" 2>/dev/null || true)"

if [ -z "${response}" ]; then
  log_err "empty response from ${RPC_URL}"
  exit 1
fi

# Consider node healthy if:
#   - it reports "result": false (fully synced), or
#   - eth_syncing structure contains "currentBlock" (sync in progress).
if echo "${response}" | grep -q '"result":false'; then
  echo "ok: nethermind is synced"
  exit 0
fi

if echo "${response}" | grep -q '"currentBlock"'; then
  echo "ok: nethermind is syncing"
  exit 0
fi

log_err "unexpected eth_syncing response: ${response}"
exit 1
