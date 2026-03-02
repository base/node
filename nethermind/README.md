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

| Setting | Default | Description |
|---------|---------|-------------|
| `NETHERMIND_TAG` | — | Version tag of Nethermind to use |
| `NETHERMIND_COMMIT` | — | Specific commit hash for verification |
| `NETHERMIND_LOG_LEVEL` | `Info` | Log verbosity level |
| `OP_SEQUENCER_HTTP` | — | Sequencer URL (required) |
| `OP_NETHERMIND_ETHSTATS_ENABLED` | — | Enable EthStats monitoring |
| `OP_NETHERMIND_ETHSTATS_NODE_NAME` | — | Node name for EthStats |
| `OP_NETHERMIND_BOOTNODES` | — | Bootnode addresses for snap sync |

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

### Health Checks

Nethermind enables health checks by default (`--HealthChecks.Enabled=true`). This exposes a health endpoint useful for monitoring and orchestration tools like Docker health checks or Kubernetes liveness probes.

### Snap Sync

For faster initial sync, enable snap sync by uncommenting the bootnode configuration:

```bash
OP_NETHERMIND_BOOTNODES=enode://...
```

> [!WARNING]
> Snap sync is experimental and may lead to syncing issues. Use with caution in production environments.

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
- Verify network bandwidth is sufficient (recommend 100+ Mbps)
- Consider using snap sync with bootnodes
- Check L1 RPC rate limits and connection stability
- Ensure NVMe storage is properly configured

#### Peer Connection Issues

If the node has difficulty finding peers:
- Ensure port 30303 (TCP/UDP) is open on your firewall
- Check that bootnodes are reachable
- Verify network connectivity from the node

#### Engine API Authentication Errors

If seeing JWT authentication errors:
- Verify `OP_NODE_L2_ENGINE_AUTH` is correctly set
- Ensure the JWT secret matches between op-node and Nethermind
- Check file permissions on the JWT secret file
