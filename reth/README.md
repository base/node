# Running a Reth Node

This is an implementation of the Reth node setup that supports Flashblocks mode based on configuration.

## Setup

- See hardware requirements mentioned in the master README
- For Flashblocks mode: Access to a Flashblocks websocket endpoint (for `RETH_FB_WEBSOCKET_URL`)
  - We provide public websocket endpoints for mainnet and devnet, included in `.env.mainnet` and `.env.sepolia`

## Node Type Selection

The node determines its mode based on the presence of the `RETH_FB_WEBSOCKET_URL` environment variable:

- **Vanilla Mode** (default): When no `RETH_FB_WEBSOCKET_URL` is provided.
- **Flashblocks Mode**: When `RETH_FB_WEBSOCKET_URL` is provided.

## Running the Node

The node follows the standard `docker-compose` workflow in the master README.

```bash
# To run Reth node with Flashblocks support, set RETH_FB_WEBSOCKET_URL in your .env file
CLIENT=reth docker-compose up
```

## Testing Flashblocks RPC Methods

When running in Flashblocks mode (with `RETH_FB_WEBSOCKET_URL` configured), you can query a pending block using the Flashblocks RPC:

```bash
curl -X POST \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["pending", false],"id":1}' \
  http://localhost:8545
```

## Additional RPC Methods

For a complete list of supported RPC methods, refer to:

- [Standard Ethereum JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Flashblocks RPC Methods](https://docs.base.org/chain/flashblocks#rpc-api) (Flashblocks mode only)
