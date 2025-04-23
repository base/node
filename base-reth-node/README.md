# Base Reth Node

Base Reth node is an extension of [Reth](https://github.com/paradigmxyz/reth) that includes various experimental features for Base.

## Features

- Full Ethereum-compatible JSON-RPC API
- Flashbots RPC methods support ([Flashblocks documentation](https://docs.base.org/chain/flashblocks#rpc-api))

## Setup

### Prerequisites

- Same requirements as the master README
- Access to a Flashblocks websocket endpoint (for `RETH_FB_WEBSOCKET_URL`)
  - We provide public websocket endpoints for mainnet and devnet, included in `.env.mainnet` and `.env.sepolia`

## Running the Node

Running the node follows the standard `docker-compose` workflow in the master README.

```bash
CLIENT=base-reth-node docker-compose up
```

## Testing Flashblocks RPC Methods

Query a pending block using the Flashblocks RPC:

```bash
curl -X POST \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["pending", false],"id":1}' \
  http://localhost:8545
```

## Additional RPC Methods

For a complete list of supported RPC methods, refer to:

- [Standard Ethereum JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Flashblocks RPC Methods](https://docs.base.org/chain/flashblocks#rpc-api)
