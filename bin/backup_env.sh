#!/usr/bin/env bash
set -euo pipefail

# Creates a timestamped backup of .env and config files.
# Backups are stored under ./backups with the current date and time.

TIMESTAMP=$(date +"%Y%m%d%H%M%S")
BACKUP_DIR="./backups/${TIMESTAMP}"
mkdir -p "${BACKUP_DIR}"

for file in .env docker-compose.yml docker-compose.override*.yml; do
  if [ -f "$file" ]; then
    cp "$file" "${BACKUP_DIR}/"
    echo "Backed up $file → ${BACKUP_DIR}/"
  fi
done

echo "Backup complete. Files are saved in ${BACKUP_DIR}"
