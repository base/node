#!/usr/bin/env bash
set -euo pipefail

RPC_URL="${1:-http://localhost:8545}"

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is not installed" >&2
  exit 1
fi

REQUEST='{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}'

echo "Checking Base node health at $RPC_URL..."
HTTP_CODE=$(curl -sS -o /tmp/base_node_health_response.json -w "%{http_code}" \
  -H "Content-Type: application/json" \
  -d "$REQUEST" \
  "$RPC_URL" || echo "000")

if [ "$HTTP_CODE" != "200" ]; then
  echo "Health check failed, HTTP status: $HTTP_CODE"
  exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
  echo "Health check succeeded (HTTP 200), but jq is not installed."
  echo "Raw response:"
  cat /tmp/base_node_health_response.json
  exit 0
fi

BLOCK_HEX=$(jq -r '.result // empty' /tmp/base_node_health_response.json 2>/dev/null || echo "")

if [ -z "$BLOCK_HEX" ] || [ "$BLOCK_HEX" = "null" ]; then
  echo "Health check failed: empty block number in response"
  cat /tmp/base_node_health_response.json
  exit 1
fi

echo "Health check OK, latest block: $BLOCK_HEX"
