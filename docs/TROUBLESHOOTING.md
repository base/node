# Base node troubleshooting FAQ

This page collects common issues seen when running a Base node and how to address them.

## The node will not start

- Verify that the RPC environment variables are set correctly (`L1_RPC_URL`, `OP_NODE_RPC_URL`).
- Check that no other process is using the configured ports (use `lsof -i :8545`).
- Look at the container logs with `docker compose logs -f` to see specific error messages.

## Sync is extremely slow

- Ensure you are using the recommended client (reth/geth/nethermind) for your hardware.
- Consider restoring from a snapshot instead of syncing from genesis.
- Check your network bandwidth and disk I/O performance.

## Logs show “unauthorized” or “invalid auth”

- Double‑check your authentication tokens or keys.
- If using Flashbots/FB endpoints, verify that `RETH_FB_WEBSOCKET_URL` is correct and reachable.

## Still stuck?

Join the community Discord’s node-operator channel or open an issue on GitHub with details about your hardware, configuration and log excerpts.
