#!/usr/bin/env bash
set -euo pipefail

OP_NODE_RPC="${OP_NODE_RPC:-http://localhost:7545}"

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required for this script."
  exit 1
fi

NOW_TS=$(date +%s)

SYNC_JSON=$(curl -sS \
  -H "Content-Type: application/json" \
  -d '{"id":1,"jsonrpc":"2.0","method":"optimism_syncStatus","params":[]}' \
  "${OP_NODE_RPC}")

L2_TS=$(echo "${SYNC_JSON}" | jq -r '.result.unsafe_l2.timestamp // 0')

if [ "${L2_TS}" -eq 0 ]; then
  echo "Could not read timestamp from response:"
  echo "${SYNC_JSON}"
  exit 1
fi

BEHIND_SECONDS=$((NOW_TS - L2_TS))
BEHIND_MINUTES=$((BEHIND_SECONDS / 60))

echo "Node is approximately ${BEHIND_MINUTES} minutes behind L2 tip."
