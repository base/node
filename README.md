# Base Node

Base is a secure, low-cost, developer-friendly Ethereum L2 built on Optimism's OP Stack. This repository contains Docker builds to run your own node on the Base network.

## Quick Start

1. Ensure you have an Ethereum L1 full node RPC available
2. Choose your network:
   - For mainnet: Use `.env.mainnet`
   - For testnet: Use `.env.sepolia`
3. Configure your L1 endpoints in the appropriate `.env` file:
   - `OP_NODE_L1_ETH_RPC=<your-preferred-l1-rpc>`
   - `OP_NODE_L1_BEACON=<your-preferred-l1-beacon>`
   - `OP_NODE_L1_BEACON_ARCHIVER=<your-preferred-l1-beacon-archiver>`

4. Start the node:
   - For mainnet (default): `docker compose up --build`
   - For testnet: `NETWORK_ENV=.env.sepolia docker compose up --build`

## Supported Clients

- **reth** (default)
- **geth**
- **nethermind**

## Requirements

### Minimum Requirements
- Modern Multicore CPU
- 32GB RAM (64GB Recommended)
- NVMe SSD drive
- Storage: (2 * current chain size + snapshot size + 20% buffer)
- Docker and Docker Compose

## Configuration

### Required Settings
- **L1 Configuration**: RPC, Beacon, and Archiver endpoints.
- **OP_NODE_L1_RPC_KIND**: Supported values include `alchemy`, `quicknode`, `infura`, `parity`, `nethermind`, `debug_geth`, `erigon`, `basic`, `any`, `standard`.

### Network Settings (Mainnet)
- `RETH_CHAIN=base`
- `OP_NODE_NETWORK=base-mainnet`
- Sequencer: `https://mainnet-sequencer.base.org`

## Snapshots
Snapshots are available to help you sync your node more quickly. See [docs.base.org](https://docs.base.org) for details.

## Supported Networks

| Network | Status |
| ------- | ------ |
| Mainnet | ✅     |
| Testnet | ✅     |

## Troubleshooting
For support please join our Discord post in #🛠｜node-operators. You can alternatively open a new GitHub issue.

## Disclaimer
THE NODE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND.
