# Troubleshooting

This document lists common issues and what information to gather before opening an issue.

## Before you open an issue
- Confirm which network you are using (mainnet vs sepolia).
- Confirm which execution client you selected (reth / geth / nethermind).
- Remove any secrets from logs (API keys, RPC tokens, private endpoints).

## Common issues

### 1) Node won’t start
Check:
- Your `.env` file selection matches your network (e.g., mainnet vs sepolia).
- Required L1 endpoints are configured (L1 RPC + beacon endpoints).
- Disk space is sufficient for the selected mode (archive vs full).

What to include in an issue:
- Which `.env` file you used
- Client name (reth/geth/nethermind)
- A short log excerpt around the failure

### 2) Slow syncing
Possible causes:
- Insufficient SSD performance or low available disk space
- Not enough RAM allocated
- L1 RPC / beacon endpoints rate-limited or unstable

What to include in an issue:
- Hardware specs (CPU/RAM/storage type)
- Client name
- Whether you use an external L1 provider (and which one)

### 3) Frequent restarts
Check:
- Resource limits (RAM, CPU)
- Host machine stability
- Configuration typos in `.env` files

What to include in an issue:
- Restart frequency
- Log excerpt around each restart
- OS + Docker version

## How to ask for help
If you open an issue, include:
- Network + client
- Environment details
- Minimal reproduction steps
- Sanitized logs (no secrets)
