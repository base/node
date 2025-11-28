#!/usr/bin/env bash
set -euo pipefail

OUTPUT_FILE="${1:-base-node-support-info.txt}"

echo "Writing support information to: ${OUTPUT_FILE}"
echo "Note: review this file for sensitive data before sharing it." > "${OUTPUT_FILE}"

section() {
  echo >> "${OUTPUT_FILE}"
  echo "===== $1 =====" >> "${OUTPUT_FILE}"
}

section "Basic system information"
{
  date
  uname -a || true
} >> "${OUTPUT_FILE}" 2>&1

section "Docker version"
{
  docker --version || true
  docker compose version || docker-compose version || true
} >> "${OUTPUT_FILE}" 2>&1

section "Docker compose status"
{
  echo "docker compose ps:"
  docker compose ps || true
} >> "${OUTPUT_FILE}" 2>&1

section "Recent logs (short tail)"
{
  echo "Note: adjust service names to match your setup if needed."
  echo

  echo "op-node (if present):"
  docker compose logs --tail=100 op-node 2>&1 || echo "op-node service not found."

  echo
  echo "reth (if present):"
  docker compose logs --tail=100 reth 2>&1 || echo "reth service not found."

  echo
  echo "geth (if present):"
  docker compose logs --tail=100 geth 2>&1 || echo "geth service not found."

  echo
  echo "nethermind (if present):"
  docker compose logs --tail=100 nethermind 2>&1 || echo "nethermind service not found."
} >> "${OUTPUT_FILE}" 2>&1

section "Environment hints"
{
  echo "List of .env-style files in the current directory:"
  ls -1 .env* 2>/dev/null || echo "No .env files found in current directory."

  echo
  echo "Note: this script does not print the contents of .env files in order"
  echo "to avoid leaking secrets. If support needs specific configuration,"
  echo "share only the relevant lines, with any secrets redacted."
} >> "${OUTPUT_FILE}" 2>&1

section "Docker images"
{
  docker compose images || true
} >> "${OUTPUT_FILE}" 2>&1

section "Disk usage (high level)"
{
  df -h . || true
} >> "${OUTPUT_FILE}" 2>&1

echo
echo "Support information collection completed."
echo "File created at: ${OUTPUT_FILE}"
echo "Please review it for secrets before sharing."
