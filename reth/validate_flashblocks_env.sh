#!/usr/bin/env bash
set -euo pipefail

# Sanity checks for a Base Reth + Flashblocks environment.
#
# This script is intended to be run before starting the reth-based node.
# It validates a minimal set of environment variables that are commonly
# required when running with Flashblocks support:
#
#   NODE_TYPE               - should be set to "base"
#   NETWORK_ENV             - points to the env file describing the network
#   RETH_FB_WEBSOCKET_URL   - websocket endpoint for Flashblocks
#
# Exit code:
#   0 - configuration looks good
#   1 - one or more validation errors were found

errors=0

log_err() {
  echo "[flashblocks-env] $*" >&2
}

require_var() {
  local name="$1"
  local value="${!name:-}"
  if [ -z "${value}" ]; then
    log_err "required environment variable ${name} is not set"
    errors=1
  fi
}

# Check that the key variables are present
require_var "NODE_TYPE"
require_var "NETWORK_ENV"
require_var "RETH_FB_WEBSOCKET_URL"

# Validate values where we can
if [ "${NODE_TYPE:-}" != "base" ]; then
  log_err "NODE_TYPE should be 'base' when running a Base Reth node (current: '${NODE_TYPE:-}')"
  errors=1
fi

fb_url="${RETH_FB_WEBSOCKET_URL:-}"
if [ -n "${fb_url}" ]; then
  case "${fb_url}" in
    ws://*|wss://*)
      # looks fine
      ;;
    *)
      log_err "RETH_FB_WEBSOCKET_URL should start with ws:// or wss:// (current: '${fb_url}')"
      errors=1
      ;;
  esac
fi

# NETWORK_ENV is expected to be a file name like ".env.mainnet" or ".env.sepolia".
# We only warn if the file does not exist in the current working directory;
# the caller may still be sourcing it from elsewhere.
if [ -n "${NETWORK_ENV:-}" ] && [ ! -f "${NETWORK_ENV}" ]; then
  log_err "NETWORK_ENV='${NETWORK_ENV}' does not point to a file in the current directory"
  errors=1
fi

if [ "${errors}" -ne 0 ]; then
  log_err "invalid Flashblocks / Base Reth environment configuration"
  exit 1
fi

echo "[flashblocks-env] configuration looks good"
exit 0
