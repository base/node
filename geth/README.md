# Running a Geth Node

This is an implementation of the op-geth (optimism-geth) node setup for running a Base node with the Geth execution client.

## Overview

[op-geth](https://github.com/ethereum-optimism/op-geth) is the Optimism implementation of the go-ethereum client, optimized for running OP Stack rollups like Base. It's a well-tested, production-ready client suitable for full node operation.

## Setup

- See hardware requirements mentioned in the master README
- Minimum 32GB RAM (64GB recommended)
- NVMe SSD storage for optimal performance

## Configuration

The Geth client is configured through environment variables in .env.mainnet or .env.sepolia. Key settings include:

### Cache Settings

Optimize cache allocation based on your available RAM:

| Setting | Default | Description |
|---------|---------|-------------|
| GETH_CACHE | 20480 | Total cache allocation in MB (20GB default) |
| GETH_CACHE_DATABASE | 20 | Percentage allocated to database cache |
| GETH_CACHE_GC | 12 | Percentage allocated to garbage collection |
| GETH_CACHE_SNAPSHOT | 24 | Percentage allocated to snapshot cache |
| GETH_CACHE_TRIE | 44 | Percentage allocated to trie cache |

### Sync Modes

| Mode | Environment Variable | Description |
|------|---------------------|-------------|
| Full Sync | OP_GETH_SYNCMODE=full | Download and verify all blocks (default) |
| Snap Sync | OP_GETH_SYNCMODE=snap | Faster initial sync (experimental) |

### Network Modes

| Mode | Environment Variable | Description |
|------|---------------------|-------------|
| Full | OP_GETH_GCMODE=full | Standard full node operation |
| Archive | OP_GETH_GCMODE=archive | Retain all historical states |

## Running the Node

The node follows the standard docker-compose workflow in the master README:

\\\ash
# Run Geth node (default when no CLIENT is specified)
docker-compose up

# Or explicitly specify Geth
CLIENT=geth docker-compose up
\\\`n
## Exposed Ports

| Port | Protocol | Description |
|------|----------|-------------|
| 8545 | HTTP | JSON-RPC endpoint |
| 8546 | WebSocket | WebSocket RPC endpoint |
| 8551 | HTTP | Engine API (authenticated) |
| 6060 | HTTP | Metrics endpoint |
| 30303 | TCP/UDP | P2P network |

## Optional Features

### EthStats Monitoring

Enable node monitoring by uncommenting in your .env file:

\\\ash
OP_GETH_ETH_STATS=nodename:secret@host:port
\\\`n
### Snap Sync (Experimental)

For faster initial sync, enable snap sync by uncommenting the bootnode configuration and setting sync mode:

\\\ash
OP_GETH_SYNCMODE=snap
OP_GETH_BOOTNODES=enode://...
\\\`n
> [!WARNING]
> Snap sync is experimental and may lead to syncing issues. Use with caution in production environments.

## Additional RPC Methods

For a complete list of supported RPC methods, refer to:

- [Standard Ethereum JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Geth RPC Documentation](https://geth.ethereum.org/docs/rpc/server)