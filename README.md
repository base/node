![Base](logo.webp)

# Base Node

Base is a secure, low-cost, developer-friendly Ethereum L2 built on Optimism's [OP Stack](https://docs.optimism.io/). This repository contains Docker builds to run your own node on the Base network.

[![Website base.org](https://img.shields.io/website-up-down-green-red/https/base.org.svg)](https://base.org)
[![Docs](https://img.shields.io/badge/docs-up-green)](https://docs.base.org/)
[![Discord](https://img.shields.io/discord/1067165013397213286?label=discord)](https://base.org/discord)
[![Twitter Base](https://img.shields.io/twitter/follow/Base?style=social)](https://x.com/Base)
[![Farcaster Base](https://img.shields.io/badge/Farcaster_Base-3d8fcc)](https://farcaster.xyz/base)

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
   NETWORK_ENV=.env.sepolia docker compose up --build

   # To use a specific client (optional):
   CLIENT=reth docker compose up --build

   # For testnet with a specific client:
   NETWORK_ENV=.env.sepolia CLIENT=reth docker compose up --build
   ```

### Supported Clients

- `reth` (default)
- `geth`
- `nethermind`

## Requirements

### Minimum Requirements

- Modern Multicore CPU
- 32GB RAM (64GB Recommended)
- NVMe SSD drive
- Storage: (2 \* [current chain size](https://base.org/stats) + [snapshot size](https://basechaindata.vercel.app) + 20% buffer) (to accommodate future growth)
- Docker and Docker Compose

### Production Hardware Specifications

The following are the hardware specifications we use in production:

#### Reth Archive Node (recommended)

- **Instance**: AWS i7i.12xlarge
- **Storage**: RAID 0 of all local NVMe drives (`/dev/nvme*`)
- **Filesystem**: ext4

#### Geth Full Node

- **Instance**: AWS i7i.12xlarge
- **Storage**: RAID 0 of all local NVMe drives (`/dev/nvme*`)
- **Filesystem**: ext4

> [!NOTE]
To run the node using a supported client, you can use the following command:
`CLIENT=supported_client docker compose up --build`
 
Supported clients:
 - reth (runs vanilla node by default, Flashblocks mode enabled by providing RETH_FB_WEBSOCKET_URL, see [Reth Node README](./reth/README.md))
 - geth
 - nethermind

## Configuration

### Required Settings

- L1 Configuration:
  - `OP_NODE_L1_ETH_RPC`: Your Ethereum L1 node RPC endpoint
  - `OP_NODE_L1_BEACON`: Your L1 beacon node endpoint
  - `OP_NODE_L1_BEACON_ARCHIVER`: Your L1 beacon archiver endpoint
  - `OP_NODE_L1_RPC_KIND`: The type of RPC provider being used (default: "debug_geth"). Supported values:
    - `alchemy`: Alchemy RPC provider
    - `quicknode`: QuickNode RPC provider
    - `infura`: Infura RPC provider
    - `parity`: Parity RPC provider
    - `nethermind`: Nethermind RPC provider
    - `debug_geth`: Debug Geth RPC provider
    - `erigon`: Erigon RPC provider
    - `basic`: Basic RPC provider (standard receipt fetching only)
    - `any`: Any available RPC method
    - `standard`: Standard RPC methods including newer optimized methods

### Network Settings

- Mainnet:
  - `RETH_CHAIN=base`
  - `OP_NODE_NETWORK=base-mainnet`
  - Sequencer: `https://mainnet-sequencer.base.org`

### Performance Settings

- Cache Settings:
  - `GETH_CACHE="20480"` (20GB)
  - `GETH_CACHE_DATABASE="20"` (4GB)
  - `GETH_CACHE_GC="12"`
  - `GETH_CACHE_SNAPSHOT="24"`
  - `GETH_CACHE_TRIE="44"`

### Optional Features

  * EthStats Monitoring (uncomment to enable)
  * Trusted RPC Mode (uncomment to enable)
  * Snap Sync (experimental)

For full configuration options, see the `.env.mainnet` file.

### Environment variables reference

The most commonly used environment variables for this repository are summarized below.  
You can set them in `.env.mainnet` / `.env.sepolia` or export them in your shell before running `docker compose`.

| Variable                      | Required | Description                                                                                  | Example                                             |
| ---------------------------- | -------- | -------------------------------------------------------------------------------------------- | --------------------------------------------------- |
| `OP_NODE_L1_ETH_RPC`         | Yes      | Ethereum L1 execution RPC endpoint used by the Base node.                                    | `https://your-l1-rpc.example`                      |
| `OP_NODE_L1_BEACON`          | Yes      | L1 beacon node endpoint.                                                                     | `https://your-l1-beacon.example`                   |
| `OP_NODE_L1_BEACON_ARCHIVER` | Yes      | L1 beacon archiver endpoint used for historical data.                                       | `https://your-l1-archiver.example`                 |
| `OP_NODE_L1_RPC_KIND`        | No       | Type of RPC provider being used. Defaults to `"debug_geth"`. See the list in **Configuration**. | `quicknode`, `alchemy`, `infura`, `nethermind`     |
| `RETH_CHAIN`                 | Yes      | Network identifier for Reth.                                                                 | `base`                                             |
| `OP_NODE_NETWORK`            | Yes      | Base network configuration name for op-node.                                                 | `base-mainnet`                                     |
| `GETH_CACHE`                 | No       | Total Geth cache size in MB.                                                                 | `20480`                                            |
| `GETH_CACHE_DATABASE`        | No       | Geth database cache (subset of `GETH_CACHE`).                                                | `20`                                               |
| `GETH_CACHE_GC`              | No       | Geth garbage collection cache percentage.                                                    | `12`                                               |
| `GETH_CACHE_SNAPSHOT`        | No       | Geth snapshot cache percentage.                                                              | `24`                                               |
| `GETH_CACHE_TRIE`            | No       | Geth trie cache percentage.                                                                  | `44`                                               |
| `NETWORK_ENV`                | No       | Path to the env file that `docker compose` should use.                                       | `.env.mainnet`, `.env.sepolia`                     |
| `CLIENT`                     | No       | L2 client to run in the Docker Compose stack.                                                | `reth`, `geth`, `nethermind`                       |

## Snapshots

Snapshots are available to help you sync your node more quickly. See [docs.base.org](https://docs.base.org/chain/run-a-base-node#snapshots) for links and more details on how to restore from a snapshot.

## Supported Networks

| Network | Status |
| ------- | ------ |
| Mainnet | ✅     |
| Testnet | ✅     |

## Troubleshooting

For support please join our [Discord](https://discord.gg/buildonbase) post in `🛠｜node-operators`. You can alternatively open a new GitHub issue.

## Disclaimer

THE NODE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. We make no guarantees about asset protection or security. Usage is subject to applicable laws and regulations.

For more information, visit [docs.base.org](https://docs.base.org/).
