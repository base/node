# Running a Nethermind Node

This is an implementation of the Nethermind node setup for running a Base node with the Nethermind execution client.

## Overview

[Nethermind](https://github.com/NethermindEth/nethermind) is a high-performance, enterprise-grade Ethereum client written in .NET. It provides an alternative execution client option for running Base nodes with excellent support for enterprise deployments.

## Setup

- See hardware requirements mentioned in the master README
- Minimum 32GB RAM (64GB recommended)
- NVMe SSD storage for optimal performance
- The Nethermind client is built using .NET SDK 9.0

## Configuration

The Nethermind client is configured through environment variables in `.env.mainnet` or `.env.sepolia`.

### Key Configuration Options

| Setting | Description |
|---------|-------------|
| `NETHERMIND_TAG` | Version tag of Nethermind to use |
| `NETHERMIND_COMMIT` | Specific commit hash for verification |
| `OP_NETHERMIND_ETHSTATS_ENABLED` | Enable EthStats monitoring |
| `OP_NETHERMIND_ETHSTATS_NODE_NAME` | Node name for EthStats |
| `OP_NETHERMIND_BOOTNODES` | Bootnode addresses for snap sync |

## Running the Node

The node follows the standard `docker-compose` workflow in the master README:

```bash
# Run Nethermind node
CLIENT=nethermind docker-compose up
```

## Architecture Support

Nethermind supports multiple architectures:

| Architecture | Identifier |
|-------------|------------|
| AMD64/x86_64 | `x64` |
| ARM64 | `arm64` |

The Dockerfile automatically detects and builds for the appropriate architecture.

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

Enable node monitoring by uncommenting in your `.env` file:

```bash
OP_NETHERMIND_ETHSTATS_ENABLED=true
OP_NETHERMIND_ETHSTATS_NODE_NAME=NethermindNode
OP_NETHERMIND_ETHSTATS_ENDPOINT=ethstats_endpoint
```

### Snap Sync

For faster initial sync, enable snap sync by uncommenting the bootnode configuration:

```bash
OP_NETHERMIND_BOOTNODES=enode://...
```

## Additional Resources

For more information about Nethermind configuration and features, refer to:

- [Nethermind Documentation](https://docs.nethermind.io/)
- [Standard Ethereum JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Nethermind GitHub Repository](https://github.com/NethermindEth/nethermind)

## Troubleshooting

### Common Issues

#### Memory Usage

Nethermind may require adjustment of .NET runtime settings for optimal memory usage. Ensure your system has sufficient RAM (64GB recommended for production).

#### Sync Performance

If experiencing slow sync performance:
- Verify network bandwidth is sufficient
- Consider using snap sync with bootnodes
- Check L1 RPC rate limits
