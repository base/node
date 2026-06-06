![Base](logo.webp)

# Base Node

Base is a secure, low-cost, developer-friendly Ethereum L2 running on the [Base stack](https://github.com/base/base). This repository contains a Docker build for running a Base node with `base-reth-node` and `base-consensus`.

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
4. Start the node:

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

### Missing L1InfoDeposit error when using pruned snapshots

If you see an error like `EngineReset(SyncStart(FromBlock(MissingL1InfoDeposit(...))))` in the consensus logs after syncing with a pruned snapshot, it means the snapshot does not contain enough historical data for the consensus client to verify the L1 deposit contract initialization.

**Solution:** Use an unpruned (archive) snapshot or a pruned snapshot that includes at least the first ~40 days of history (approximately block 45,000,000 as of mid‑2026). Alternatively, sync from scratch without a snapshot, or use an archive node for the execution client.

**Steps to resolve:**

1. Verify the snapshot age: check the block number associated with the snapshot (provided in the snapshot name or description).

2. If the snapshot is older than ~40 days, download a newer snapshot from the official snapshots page: https://docs.base.org/chain/run-a-base-node#snapshots

3. If you prefer not to use snapshots, remove the snapshot data directory and let the node sync from genesis (this will take longer but guarantees completeness).

4. Ensure that both execution and consensus clients are using the same snapshot data directory (they share `/data` in the default compose setup).

```
# Example: remove old data and restart
rm -rf ./reth-data/*
docker compose up --build
```

**Note:** This issue is not a bug in the node software but a limitation of pruned snapshots that discard historic state needed for deposit contract verification.

## Disclaimer

THE NODE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. We make no guarantees about asset protection or security. Usage is subject to applicable laws and regulations.

For more information, visit [docs.base.org](https://docs.base.org/).
