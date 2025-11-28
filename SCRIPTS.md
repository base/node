# Helper scripts

This repository can include optional helper scripts to operate a Base node.

## scripts/check-sync-status.sh

This script queries the `optimism_syncStatus` method on the op-node RPC endpoint
and prints how many minutes the node is behind the L2 tip.

Usage example:

```bash
OP_NODE_RPC=http://localhost:7545 ./scripts/check-sync-status.sh
