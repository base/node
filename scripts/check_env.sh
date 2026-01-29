#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="${1:-.env.mainnet}"

if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: Env file not found: $ENV_FILE"
  exit 1
fi

required_keys=("HOST_DATA_DIR")

missing=0
for key in "${required_keys[@]}"; do
  if ! grep -qE "^${key}=" "$ENV_FILE"; then
    echo "Missing required key: ${key}"
    missing=1
  fi
done

if [ "$missing" -ne 0 ]; then
  exit 1
fi

echo "Env check passed."
