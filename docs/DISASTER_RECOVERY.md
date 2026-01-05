# Disaster recovery guide

This guide explains how to back up and restore a Base node in the event of data loss.

## Backups

- Regularly compress your data directory (`/data/reth`, `/data/geth`, etc.) and store it offsite.
- Back up your `.env` file and any JWT secrets.

## Restore from snapshot

1. Stop the node and remove the current data directory.
2. Extract your backup archive to the original data path.
3. Restart the node and verify it syncs from the restored state.

## Snapshots

For faster recovery, use provider snapshots when available. Check official Base docs for instructions.
