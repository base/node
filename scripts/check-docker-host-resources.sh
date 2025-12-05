#!/usr/bin/env bash
set -euo pipefail

echo "=== Base node Docker host resources check ==="
echo

# 1. Docker availability
echo "[1/4] Checking Docker..."
if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is not installed or not on PATH."
  echo "Please install Docker and try again."
  echo
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  echo "Docker is installed but 'docker info' failed."
  echo "Make sure the Docker daemon is running and you have permission to use it."
  echo
  exit 1
fi

echo "Docker is available and the daemon is responding."
echo

# 2. Memory information (from docker info)
echo "[2/4] Checking memory configuration..."
TOTAL_MEM_BYTES="$(docker info --format '{{.MemTotal}}' 2>/dev/null || echo 0)"

if [ "${TOTAL_MEM_BYTES}" -gt 0 ] 2>/dev/null; then
  TOTAL_MEM_GIB=$(( TOTAL_MEM_BYTES / 1024 / 1024 / 1024 ))
  echo "Total memory reported by Docker: ${TOTAL_MEM_GIB} GiB"
else
  echo "Could not determine memory from 'docker info'."
  TOTAL_MEM_GIB=0
fi

if [ "${TOTAL_MEM_GIB}" -lt 32 ] 2>/dev/null; then
  echo "Warning: less than the recommended 32 GiB of RAM. Base nodes are resource-intensive."
  echo "You may experience slow sync or instability on lower-memory hosts."
else
  echo "Memory appears to meet or exceed the recommended minimum of 32 GiB."
fi

echo

# 3. CPU information (from docker info)
echo "[3/4] Checking CPU configuration..."
CPU_COUNT="$(docker info --format '{{.NCPU}}' 2>/dev/null || echo 0)"

if [ "${CPU_COUNT}" -gt 0 ] 2>/dev/null; then
  echo "Docker reports ${CPU_COUNT} CPUs available to containers."
  if [ "${CPU_COUNT}" -lt 4 ] 2>/dev/null; then
    echo "Warning: fewer than 4 CPUs may lead to slower performance."
  fi
else
  echo "Could not determine CPU count from 'docker info'."
fi

echo

# 4. Storage driver and notes
echo "[4/4] Checking storage driver..."
STORAGE_DRIVER="$(docker info --format '{{.Driver}}' 2>/dev/null || echo "")"

if [ -n "${STORAGE_DRIVER}" ]; then
  echo "Docker storage driver: ${STORAGE_DRIVER}"
else
  echo "Could not determine storage driver from 'docker info'."
fi

echo
echo "Note:"
echo "- For production-like setups, NVMe SSD storage is strongly recommended for Base nodes."
echo "- Make sure your data directories for the client are placed on fast disks."
echo
echo "Docker host resources check complete."
echo "Use these values as a sanity check against the hardware recommendations in the README."
