#!/usr/bin/env bash
set -euo pipefail

SERVICE="${1:-}"

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is not installed or not in PATH" >&2
  exit 1
fi

if [ -n "$SERVICE" ]; then
  echo "Tailing logs for service: $SERVICE"
  docker compose logs -f "$SERVICE"
else
  echo "Tailing logs for all services..."
  docker compose logs -f
fi
