#!/usr/bin/env bash
set -euo pipefail

# Small helper script to show the status of Base node Docker services.

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yml}"

echo "Using compose file: ${COMPOSE_FILE}"
echo

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is not installed or not on PATH."
  exit 1
fi

if ! command -v docker compose >/dev/null 2>&1 && ! command -v docker-compose >/dev/null 2>&1; then
  echo "docker compose (or docker-compose) is not available."
  exit 1
fi

if [ -f "${COMPOSE_FILE}" ]; then
  echo "Listing services from ${COMPOSE_FILE}:"
else
  echo "Warning: ${COMPOSE_FILE} not found in the current directory."
fi

if command -v docker compose >/dev/null 2>&1; then
  docker compose ps || true
else
  docker-compose ps || true
fi
