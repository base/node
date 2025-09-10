![Base](logo.webp)

# Base Node

Base is a secure, low-cost, developer-friendly Ethereum L2 built on Optimism's [OP Stack](https://stack.optimism.io/). This repository contains Docker builds to run your own node on the Base network.

[![Website base.org](https://img.shields.io/website-up-down-green-red/https/base.org.svg)](https://base.org)
[![Docs](https://img.shields.io/badge/docs-up-green)](https://docs.base.org/)
[![Discord](https://img.shields.io/discord/1067165013397213286?label=discord)](https://base.org/discord)
[![Twitter Base](https://img.shields.io/twitter/follow/Base?style=social)](https://x.com/Base)
[![Farcaster Base](https://img.shields.io/badge/Farcaster_Base-3d8fcc)](https://farcaster.xyz/base)

## Quick Start

1. **Prerequisites**: Ensure you have an Ethereum L1 full node RPC available
2. **Choose your network**:
   - For mainnet: Use `.env.mainnet`
   - For testnet: Use `.env.sepolia`
3. **Set up data directory**:
   ```bash
   # Create a directory for blockchain data (adjust path as needed)
   export HOST_DATA_DIR=./data
   mkdir -p $HOST_DATA_DIR
   ```
4. **Configure your L1 endpoints** in the appropriate `.env` file:
   ```bash
   OP_NODE_L1_ETH_RPC=<your-preferred-l1-rpc>
   OP_NODE_L1_BEACON=<your-preferred-l1-beacon>
   OP_NODE_L1_BEACON_ARCHIVER=<your-preferred-l1-beacon-archiver>
   ```
5. **Start the node**:

   ```bash
   # For mainnet (default):
   docker compose up --build

   # For testnet:
   NETWORK_ENV=.env.sepolia docker compose up --build

   # To use a specific client (optional):
   CLIENT=reth docker compose up --build

   # For testnet with a specific client:
   NETWORK_ENV=.env.sepolia CLIENT=reth docker compose up --build

   # Run in detached mode (background):
   docker compose up --build -d

   # View logs:
   docker compose logs -f

   # Stop the node:
   docker compose down
   ```

### Supported Clients

- `geth` (default)
- `reth`
- `nethermind`

## Requirements

### Minimum Requirements

- Modern Multicore CPU
- 32GB RAM (64GB Recommended)
- NVMe SSD drive
- Storage: (2 \* [current chain size](https://base.org/stats) + [snapshot size](https://basechaindata.vercel.app) + 20% buffer (to accommodate future growth)
- Docker and Docker Compose

### Production Hardware Specifications

The following are the hardware specifications we use in production:

#### Geth Full Node

- **Instance**: AWS i4i.12xlarge
- **Storage**: RAID 0 of all local NVMe drives (`/dev/nvme*`)
- **Filesystem**: ext4

#### Reth Archive Node

- **Instance**: AWS i7ie.6xlarge
- **Storage**: RAID 0 of all local NVMe drives (`/dev/nvme*`)
- **Filesystem**: ext4

[!NOTE]
To run the node using a supported client, you can use the following command:
`CLIENT=supported_client docker compose up --build`
 
Supported clients:
 - geth
 - reth (with Flashblocks support option, see [Reth Node README](./reth/README.md))
 - nethermind

## Configuration

### Required Environment Variables

#### Data Directory
- `HOST_DATA_DIR`: Host directory path for blockchain data storage (default: `./data`)
  ```bash
  export HOST_DATA_DIR=/path/to/your/data/directory
  ```

#### L1 Configuration
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

#### Client Selection
- `CLIENT`: Execution client to use (`geth`, `reth`, or `nethermind`)
- `NODE_TYPE`: For Reth client, choose `vanilla` (default) or `base` (with Flashblocks support)

#### JWT Authentication
- `OP_NODE_L2_ENGINE_AUTH_RAW`: JWT secret for engine API authentication (required)
- `OP_NODE_L2_ENGINE_AUTH`: Path to JWT secret file (defined in .env files)

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

- EthStats Monitoring (uncomment to enable)
- Trusted RPC Mode (uncomment to enable)
- Snap Sync (experimental)

For full configuration options, see the `.env.mainnet` file.

## Snapshots

Snapshots are available to help you sync your node more quickly. See [docs.base.org](https://docs.base.org/chain/run-a-base-node#snapshots) for links and more details on how to restore from a snapshot.

## Deployment Examples

### Development Setup
```bash
# Quick development setup with minimal resources
export HOST_DATA_DIR=./dev-data
mkdir -p $HOST_DATA_DIR

# Use testnet for development
NETWORK_ENV=.env.sepolia CLIENT=geth docker compose up --build
```

### Production Setup
```bash
# Production setup with optimized settings
export HOST_DATA_DIR=/opt/base-node/data
mkdir -p $HOST_DATA_DIR

# Set production environment variables
export GETH_CACHE="32768"      # 32GB cache
export GETH_CACHE_DATABASE="40" # 8GB database cache

# Run with specific client
CLIENT=reth docker compose up --build -d
```

### Multi-Client Testing
```bash
# Test different clients
CLIENT=geth docker compose up --build
CLIENT=reth docker compose up --build  
CLIENT=nethermind docker compose up --build
```

### Reth with Flashblocks Support
```bash
# Enable Base-specific features
NODE_TYPE=base CLIENT=reth docker compose up --build
```

### Monitoring and Logs
```bash
# Run with log monitoring
docker compose up --build -d
docker compose logs -f --tail=100

# Monitor specific service
docker compose logs -f execution
docker compose logs -f node
```

## Supported Networks

| Network | Status |
| ------- | ------ |
| Mainnet | ✅     |
| Testnet | ✅     |

## Troubleshooting

### Common Docker Issues

#### 1. Port Already in Use
```bash
# Error: Port 8545 is already in use
# Solution: Stop conflicting services or change ports
docker compose down
sudo lsof -i :8545  # Check what's using the port
```

#### 2. Permission Denied on Data Directory
```bash
# Error: Permission denied when accessing data directory
# Solution: Fix directory permissions
sudo chown -R $USER:$USER ./data
chmod -R 755 ./data
```

#### 3. Out of Disk Space
```bash
# Error: No space left on device
# Solution: Check disk usage and clean up
df -h
docker system prune -a  # Clean up Docker images
```

#### 4. Memory Issues
```bash
# Error: Container killed due to OOM
# Solution: Increase Docker memory limits or reduce cache settings
# Edit docker-compose.yml to add memory limits:
# deploy:
#   resources:
#     limits:
#       memory: 32G
```

#### 5. L1 RPC Connection Issues
```bash
# Error: Failed to connect to L1 RPC
# Solution: Verify your L1 endpoint and credentials
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  $OP_NODE_L1_ETH_RPC
```

### Docker Compose Commands

#### Basic Operations
```bash
# Start services
docker compose up --build

# Start in background
docker compose up --build -d

# View logs
docker compose logs -f
docker compose logs -f execution  # Specific service logs

# Stop services
docker compose down

# Restart services
docker compose restart

# Rebuild and restart
docker compose up --build --force-recreate
```

#### Debugging
```bash
# Check service status
docker compose ps

# Execute commands in running container
docker compose exec execution bash
docker compose exec node bash

# Check container resource usage
docker stats

# View detailed container info
docker compose config
```

### Performance Optimization

#### For Low-Memory Systems (16GB RAM)
```bash
# Reduce cache settings in your .env file:
GETH_CACHE="8192"        # 8GB instead of 20GB
GETH_CACHE_DATABASE="10"  # 2GB instead of 4GB
GETH_CACHE_GC="6"         # Reduce GC frequency
```

#### For High-Performance Systems
```bash
# Increase cache settings:
GETH_CACHE="32768"       # 32GB cache
GETH_CACHE_DATABASE="40" # 8GB database cache
GETH_CACHE_TRIE="64"     # Larger trie cache
```

### Getting Help

For additional support:
- Join our [Discord](https://discord.gg/buildonbase) and post in `🛠｜node-operators`
- Open a new [GitHub issue](https://github.com/base/node/issues)
- Check the [Base documentation](https://docs.base.org/)

## Disclaimer

THE NODE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND. We make no guarantees about asset protection or security. Usage is subject to applicable laws and regulations.

For more information, visit [docs.base.org](https://docs.base.org/).
