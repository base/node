![Base](logo.webp)

# Base Node

Base is a secure, low-cost, developer-friendly Ethereum L2 built on Optimism's [OP Stack](https://docs.optimism.io/). This repository contains a Docker build for running a Base node with `base-reth-node` and `base-consensus`.

[![Website base.org](https://img.shields.io/website-up-down-green-red/https/base.org.svg)](https://base.org)
[![Docs](https://img.shields.io/badge/docs-up-green)](https://docs.base.org/)
[![Discord](https://img.shields.io/discord/1067165013397213286?label=discord)](https://base.org/discord)
[![Twitter Base](https://img.shields.io/twitter/follow/Base?style=social)](https://x.com/Base)
[![Farcaster Base](https://img.shields.io/badge/Farcaster_Base-3d8fcc)](https://farcaster.xyz/base)

## Prerequisites

Before running a Base node, ensure you have the following installed and configured:

- **Docker** v20.10 or higher — [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose** v2.0 or higher — [Install Docker Compose](https://docs.docker.com/compose/install/)
- **An Ethereum L1 full node RPC endpoint** — options include [Alchemy](https://www.alchemy.com/), [Infura](https://www.infura.io/), [QuickNode](https://www.quicknode.com/), or your own self-hosted node
- **An Ethereum L1 Beacon endpoint** — required for consensus layer communication

To verify your Docker installation:

```bash
docker --version
docker compose version
```

## Quick Start

1. Ensure you have an Ethereum L1 full node RPC and beacon endpoint available.
2. Choose your network:
   - For mainnet: use `.env.mainnet`
   - For testnet: use `.env.sepolia`
3. Configure your L1 endpoints in the appropriate `.env` file:
   ```bash
   BASE_NODE_L1_ETH_RPC=<your-preferred-l1-rpc>
   BASE_NODE_L1_BEACON=<your-preferred-l1-beacon>
   ```
4. Start the node:

   ```bash
   # For mainnet (default):
   docker compose up --build

   # For testnet:
   NETWORK_ENV=.env.sepolia docker compose up --build
   ```

5. Verify your node is running correctly:

   ```bash
   # Check that containers are running
   docker compose ps

   # Check node logs for sync progress
   docker compose logs -f

   # Query the node's current block number
   curl -X POST \
     --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
     http://localhost:8545
   ```

   A successful response will look like:
   ```json
   {"jsonrpc":"2.0","id":1,"result":"0x..."}
   ```

   It may take several minutes for the node to begin syncing. Use snapshots (see below) to speed up initial sync significantly.

## Supported Clients

- Execution: `base-reth-node`
- Consensus: `base-consensus`

## Requirements

### Minimum Requirements

- Modern multicore CPU
- 32GB RAM (64GB recommended)
- NVMe SSD drive
- Storage: `(2 × current chain size) + snapshot size + 20% buffer` to accommodate future growth
  - Example: if the current chain size is 500GB and snapshot is 200GB, you need at least `(2 × 500) + 200 + 20% = ~1.44TB`
  - Check current chain size at [base.org/stats](https://base.org/stats) and snapshot size at [basechaindata.vercel.app](https://basechaindata.vercel.app)
- Docker and Docker Compose (see [Prerequisites](#prerequisites))

### Production Hardware Specifications

The following are the hardware specifications we use in production:

#### Reth Archive Node (recommended)

- **Instance**: AWS i7i.12xlarge
- **Storage**: RAID 0 of all local NVMe drives (`/dev/nvme*`)
- **Filesystem**: ext4

## Configuration

### Required Settings

- `BASE_NODE_L1_ETH_RPC`: your Ethereum L1 node RPC endpoint
- `BASE_NODE_L1_BEACON`: your L1 beacon node endpoint
- `BASE_NODE_NETWORK`: `base` or `base-sepolia`
- `RETH_CHAIN`: `base` or `base-sepolia`

### Network Settings

- Mainnet:
  - `RETH_CHAIN=base`
  - `BASE_NODE_NETWORK=base`
  - Sequencer: `https://mainnet-sequencer.base.org`
- Sepolia:
  - `RETH_CHAIN=base-sepolia`
  - `BASE_NODE_NETWORK=base-sepolia`
  - Sequencer: `https://sepolia-sequencer.base.org`

### Optional Features

- **Flashblocks**: set `RETH_FB_WEBSOCKET_URL`. When set, the execution client runs in Flashblocks mode; otherwise it runs in vanilla mode.
- **Follow mode**: set `BASE_NODE_SOURCE_L2_RPC`
- **Pruning**: set `RETH_PRUNING_ARGS`
- **L1 Verifier Confirmation Depth**: set `BASE_NODE_VERIFIER_L1_CONFS` to configure the number of L1 confirmations the verifier waits before processing. Can also be set via the `--l1.verifier-confs` CLI flag. Useful for node operators who want to adjust confirmation depth based on their L1 provider's reliability.

For full configuration options, see `.env.mainnet` or `.env.sepolia`.

### Testing Flashblocks RPC Methods

When running in Flashblocks mode, you can query a pending block using the Flashblocks RPC:

```bash
curl -X POST \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["pending", false],"id":1}' \
  http://localhost:8545
```

## Snapshots

Snapshots are available to help you sync your node more quickly. See [docs.base.org](https://docs.base.org/chain/run-a-base-node#snapshots) for links and more details on how to restore from a snapshot.

## Supported Networks

| Network | Status |
| ------- | ------ |
| Mainnet | ✅ |
| Testnet | ✅ |

## Upgrading

When a new version of `base-reth-node` or `base-consensus` is released, update your node by pulling the latest images and restarting:

```bash
docker compose pull
docker compose down
docker compose up --build
```

Check the [releases page](https://github.com/base-org/node/releases) for changelog and any breaking changes before upgrading.

> **Note for pruned node operators**: some releases may be mandatory for pruned nodes while remaining optional for archive nodes. Always check the release notes to determine whether an upgrade is required for your node type.

## Troubleshooting

### Common Issues

**Node not syncing:**
- Verify your L1 RPC and beacon endpoints are accessible and not rate-limited
- Check logs with `docker compose logs -f` for error messages
- Ensure your storage has sufficient free space

**Container exits immediately:**
- Run `docker compose logs` to see the error output
- Verify all required environment variables are set correctly in your `.env` file

**RPC not responding:**
- Confirm the node is fully started with `docker compose ps`
- Check that port `8545` is not blocked by a firewall

For additional support, join our [Discord](https://discord.gg/buildonbase) and post in `🛠｜node-operators`. You can alternatively open a new GitHub issue.

## Disclaimer

THE NODE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. We make no guarantees about asset protection or security. Usage is subject to applicable laws and regulations.

For more information, visit [docs.base.org](https://docs.base.org/).