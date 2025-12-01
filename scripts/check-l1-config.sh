#!/usr/bin/env bash
set -euo pipefail

echo "=== Base node L1 configuration check ==="
echo

ENV_MAINNET=".env.mainnet"
ENV_SEPOLIA=".env.sepolia"

# Pick which env file to inspect.
# Default: mainnet; override with NETWORK_ENV_FILE if desired.
ENV_FILE="${NETWORK_ENV_FILE:-${ENV_MAINNET}}"

if [ ! -f "${ENV_FILE}" ]; then
  echo "Environment file ${ENV_FILE} not found."
  echo "Make sure you have created and configured ${ENV_MAINNET} or ${ENV_SEPOLIA}."
  echo "You can switch the file by setting NETWORK_ENV_FILE before running this script."
  echo
  exit 1
fi

echo "[1/3] Using environment file: ${ENV_FILE}"
echo

required_var() {
  local var_name="$1"
  local value

  # Extract variable value from the env file (ignoring comments).
  value="$(grep -E "^${var_name}=" "${ENV_FILE}" | head -n 1 | cut -d'=' -f2- || true)"

  if [ -z "${value}" ]; then
    echo "  - ${var_name}: MISSING or EMPTY"
  else
    echo "  - ${var_name}: set (value hidden)"
  fi
}

optional_var() {
  local var_name="$1"
  local value

  value="$(grep -E "^${var_name}=" "${ENV_FILE}" | head -n 1 | cut -d'=' -f2- || true)"

  if [ -z "${value}" ]; then
    echo "  - ${var_name}: not set (using default behaviour)"
  else
    echo "  - ${var_name}: set (value hidden)"
  fi
}

echo "[2/3] Required L1 settings:"
required_var "OP_NODE_L1_ETH_RPC"
required_var "OP_NODE_L1_BEACON"
required_var "OP_NODE_L1_BEACON_ARCHIVER"
echo

echo "[3/3] Optional L1 settings:"
optional_var "OP_NODE_L1_RPC_KIND"
optional_var "OP_NODE_L1_TRUST_RPC"
echo

echo "Summary:"
echo "- Ensure all required variables above are set to reachable, fully synced L1 endpoints."
echo "- OP_NODE_L1_RPC_KIND controls how the node interacts with the provider."
echo "- OP_NODE_L1_TRUST_RPC should only be enabled if you trust the RPC source."
echo
echo "If any required variable is reported as missing or empty, please update ${ENV_FILE}"
echo "before starting the node with docker compose."
