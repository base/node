#!/usr/bin/env bash
set -euo pipefail

# Compresses all log files into a tar.gz archive with a timestamp and clears the original files.
#
# Usage:
#   ./bin/archive_logs.sh
#
# Archived files are stored in ./archives/ with a filename like logs-YYYYMMDDHHMM.tar.gz

LOG_DIR="./logs"
ARCHIVE_DIR="./archives"
TIMESTAMP=$(date +"%Y%m%d%H%M")

mkdir -p "${ARCHIVE_DIR}"

ARCHIVE_NAME="${ARCHIVE_DIR}/logs-${TIMESTAMP}.tar.gz"
tar -czf "${ARCHIVE_NAME}" -C "${LOG_DIR}" .
echo "Logs archived to ${ARCHIVE_NAME}"

# Remove original logs
find "${LOG_DIR}" -type f -name "*.log" -delete
echo "Original log files removed."
