![Base](logo.webp)

# Base Node

Base is a secure, low-cost, developer-friendly Ethereum L2 built on Optimism's [OP Stack](https://docs.optimism.io/). This repository contains a Docker build for running a Base node with `base-reth-node` and `base-consensus`.

[![Website base.org](https://img.shields.io/website-up-down-green-red/https/base.org.svg)](https://base.org)
[![Docs](https://img.shields.io/badge/docs-up-green)](https://docs.base.org/)
[![Discord](https://img.shields.io/discord/1067165013397213286?label=discord)](https://base.org/discord)
[![Twitter Base](https://img.shields.io/twitter/follow/Base?style=social)](https://x.com/Base)
[![Farcaster Base](https://img.shields.io/badge/Farcaster_Base-3d8fcc)](https://farcaster.xyz/base)

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
4. Set the data directory (optional but recommended):
   ```bash
   HOST_DATA_DIR=./reth-data
   ```
   This path on the host machine is where chain data will be persisted. If not set, Docker will default to `./reth-data`.
5. Start the node:

   ```bash
   # For mainnet (default):
   docker compose up --build

   # For testnet:
   NETWORK_ENV=.env.sepolia docker compose up --build
   ```

## Supported Clients

- Execution: `base-reth-node`
- Consensus: `base-consensus`

## Requirements

### Minimum Requirements

- Modern multicore CPU
- 32GB RAM (64GB recommended)
- NVMe SSD drive
- Storage: (2 * [current chain size](https://base.org/stats) + [snapshot size](https://basechaindata.vercel.app) + 20% buffer) to accommodate future growth
- Docker and Docker Compose

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

- Flashblocks: set `RETH_FB_WEBSOCKET_URL`. When set, the execution client runs in Flashblocks mode; otherwise it runs in vanilla mode.
- Follow mode: set `BASE_NODE_SOURCE_L2_RPC`
- Pruning: set `RETH_PRUNING_ARGS`

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

## Troubleshooting

For support please join our [Discord](https://discord.gg/buildonbase) and post in `🛠｜node-operators`. You can alternatively open a new GitHub issue.

## Disclaimer

THE NODE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. We make no guarantees about asset protection or security. Usage is subject to applicable laws and regulations.

For more information, visit [docs.base.org](https://docs.base.org/).
