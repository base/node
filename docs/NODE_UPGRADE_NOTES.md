# Node Upgrade Notes

This document summarizes a few best practices when upgrading a Base node.

## 1. Back up data

Before upgrading, back up:

- the data directory (for example `geth-data`),
- any local configuration files (`.env`, `config/*.env`),
- Docker compose files if you have customized them.

## 2. Check release notes

Always check the release notes of the new version for:

- breaking changes,
- required flags or configuration changes,
- migration steps.

## 3. Rolling upgrade (Docker)

For Docker-based setups:

1. Stop the existing containers:

       docker compose down

2. Pull the new images:

       docker pull <new-image-tags>

3. Start the node again:

       docker compose up -d

4. Monitor logs for a few minutes to ensure the node syncs correctly.

## 4. Rollback plan

Keep a short rollback plan:

- how to restore the previous image tag,
- how to restore the backed up data directory.

Having this written down reduces downtime if something goes wrong.
