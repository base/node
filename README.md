![Base](logo.webp)

# Base Node

Base is a secure, low-cost, developer-friendly Ethereum L2 built on Optimism's [OP Stack](https://docs.optimism.io/). This repository contains Docker builds to run your own node on the Base network.

[![Website base.org](https://img.shields.io/website-up-down-green-red/https/base.org.svg)](https://base.org)
[![Docs](https://img.shields.io/badge/docs-up-green)](https://docs.base.org/)
[![Discord](https://img.shields.io/discord/1067165013397213286?label=discord)](https://base.org/discord)
[![Twitter Base](https://img.shields.io/twitter/follow/Base?style=social)](https://x.com/Base)
[![Farcaster Base](https://img.shields.io/badge/Farcaster_Base-3d8fcc)](https://farcaster.xyz/base)

## Quick Start

1. Ensure you have an Ethereum L1 full node RPC available.
2. Choose your network:
   - For mainnet: Use `.env.mainnet`
   - For testnet: Use `.env.sepolia`
3. Configure your L1 endpoints in the appropriate `.env` file:
   ```bash
   OP_NODE_L1_ETH_RPC=<your-preferred-l1-rpc>
   OP_NODE_L1_BEACON=<your-preferred-l1-beacon>
   OP_NODE_L1_BEACON_ARCHIVER=<your-preferred-l1-beacon-archiver>

  4. Start the node:
   ```bash
   # For mainnet (default):
   docker compose up --build

   # For testnet:
   NETWORK_ENV=.env.sepolia docker compose up --build

 ## Requirements

### Minimum Requirements
- Modern Multicore CPU
- 32GB RAM (64GB Recommended)
- NVMe SSD drive
- Storage: (2 * [current chain size](https://base.org/stats) + [snapshot size](https://basechaindata.vercel.app) + 20% buffer)
- Docker and Docker Compose

## Supported Networks

| Network | Status |
| ------- | ------ |
| Mainnet | ✅     |
| Testnet | ✅     |

## Troubleshooting

For support please join our [Discord](https://discord.gg/buildonbase) post in `🛠｜node-operators`. You can alternatively open a new GitHub issue.

## Disclaimer

THE NODE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. Usage is subject to applicable laws and regulations. For more information, visit [docs.base.org](https://docs.base.org/).
