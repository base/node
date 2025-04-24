![Base](logo.webp)

# Base Node

Base is a secure, low-cost, developer-friendly Ethereum L2 built on Optimism's [OP Stack](https://stack.optimism.io/). This repository contains Docker builds to run your own node on the Base network.

[![Website base.org](https://img.shields.io/website-up-down-green-red/https/base.org.svg)](https://base.org)
[![Docs](https://img.shields.io/badge/docs-up-green)](https://docs.base.org/)
[![Discord](https://img.shields.io/discord/1067165013397213286?label=discord)](https://base.org/discord)
[![Twitter Base](https://img.shields.io/twitter/follow/Base?style=social)](https://x.com/Base)

## Quick Start

1. Ensure you have an Ethereum L1 full node RPC available
2. Choose your network:
   - For mainnet: Use `.env.mainnet`
   - For testnet: Use `.env.sepolia`
3. Configure your L1 endpoints in the appropriate `.env` file:
   ```bash
   OP_NODE_L1_ETH_RPC=<your-preferred-l1-rpc>
   OP_NODE_L1_BEACON=<your-preferred-l1-beacon>
   OP_NODE_L1_BEACON_ARCHIVER=<your-preferred-l1-beacon-archiver>
   ```
4. Start the node:

   ```bash
   # For mainnet (default):
   docker compose up --build

   # For testnet:
   cp .env.sepolia .env.mainnet
   docker compose up --build

   # To use a specific client (optional):
   CLIENT=reth docker compose up --build
   ```

### Supported Clients

- `geth` (default)
- `reth`
- `nethermind`

## Requirements

- Modern multi-core CPU
- 16 GB RAM (32 GB recommended)
- NVMe SSD drive
- Storage: (2 \* current_chain_size) + snapshot_size + 20% buffer
- Docker and Docker Compose

## Configuration

### Required Settings

- L1 Configuration:
  - `OP_NODE_L1_ETH_RPC`: Your Ethereum L1 node RPC endpoint
  - `OP_NODE_L1_BEACON`: Your L1 beacon node endpoint
  - `OP_NODE_L1_BEACON_ARCHIVER`: Your L1 beacon archiver endpoint

### Network Settings

- Mainnet:
  - `RETH_CHAIN=base`
  - `OP_NODE_NETWORK=base-mainnet`
  - Sequencer: `https://mainnet-sequencer.base.org`

### Performance Settings

- Cache Settings:
  - `GETH_CACHE=51200`
  - `GETH_CACHE_DATABASE=8`
  - `GETH_CACHE_GC=0`
  - `GETH_CACHE_SNAPSHOT=36`
  - `GETH_CACHE_TRIE=56`

### Optional Features

- EthStats Monitoring (uncomment to enable)
- Trusted RPC Mode (uncomment to enable)
- Snap Sync (experimental)

For full configuration options, see the `.env.mainnet` file.

## Supported Networks

| Network | Status |
| ------- | ------ |
| Mainnet | ✅     |
| Testnet | ✅     |

## Troubleshooting

For support:

1. Join our [Discord](https://discord.gg/buildonbase)
2. Connect your GitHub account via `server menu` > `Linked Roles`
3. Post in `#🛟|developer-support` or `🛠｜node-operators`

## Disclaimer

THE NODE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. We make no guarantees about asset protection or security. Usage is subject to applicable laws and regulations.

For more information, visit [docs.base.org](https://docs.base.org/).
