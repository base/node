#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="${1:-.env.mainnet}"

REQUIRED_VARS=(
  "OP_NODE_L1_ETH_RPC"
  "OP_NODE_L1_BEACON"
  "OP_NODE_L1_BEACON_ARCHIVER"
  "OP_NODE_L1_RPC_KIND"
)

PLACEHOLDER_PATTERNS=(
  "<your-preferred-l1-rpc>"
  "<your-preferred-l1-beacon>"
  "<your-preferred-l1-beacon-archiver>"
)

echo "Checking required L1 configuration variables in: ${ENV_FILE}"
echo

if [ ! -f "${ENV_FILE}" ]; then
  echo "Error: file '${ENV_FILE}' does not exist."
  echo "Pass a different file path as the first argument if needed."
  exit 1
fi

missing=0
placeholder=0

get_var_line() {
  local name="$1"
  # Take the last assignment in the file if there are multiple
  grep -E "^[[:space:]]*${name}=" "${ENV_FILE}" | tail -n 1 || true
}

strip_value() {
  sed 's/^[[:space:]]*[A-Za-z0-9_]\{1,\}=[[:space:]]*//' | tr -d '"' | tr -d "'"
}

for var in "${REQUIRED_VARS[@]}"; do
  line="$(get_var_line "${var}")"
  if [ -z "${line}" ]; then
    echo "Missing variable: ${var}"
    missing=$((missing + 1))
    continue
  fi

  value="$(printf '%s\n' "${line}" | strip_value)"

  if [ -z "${value}" ]; then
    echo "Empty value for: ${var}"
    missing=$((missing + 1))
    continue
  fi

  for pat in "${PLACEHOLDER_PATTERNS[@]}"; do
    if printf '%s\n' "${value}" | grep -q "${pat}"; then
      echo "Placeholder value detected for: ${var}"
      placeholder=$((placeholder + 1))
      break
    fi
  done
done

echo
if [ "${missing}" -eq 0 ] && [ "${placeholder}" -eq 0 ]; then
  echo "OK: all required variables are present and non-empty."
  exit 0
fi

echo "Summary:"
echo "  Missing or empty variables    : ${missing}"
echo "  Variables with placeholder(s) : ${placeholder}"
echo
echo "Please update '${ENV_FILE}' with valid endpoints before running the node."
exit 1
