# Node Healthcheck Runbook

This runbook provides quick checks to validate your node is healthy.

1) Containers are up

Run:
  docker compose ps

You should see expected services in running state.

2) Logs show progress (no loops)

Run:
  docker compose logs --since=10m --tail=200

Look for:
- steady block processing or syncing messages
- no repeated crash or restart loops

3) RPC answers

Check block number from the execution endpoint by running:
  curl -s -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545

A valid hex block number means the RPC endpoint is responsive.

4) Disk space

Run:
  df -h
  du -sh ./data 2>/dev/null || true

Ensure there is enough free disk space and the data directory is not growing unexpectedly fast.

5) Common quick fixes

Restart services:
  docker compose restart

Pull latest images and restart:
  docker compose pull
  docker compose up -d

If the node state looks corrupted, restore from a known-good snapshot following the official Base documentation.
