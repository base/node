# Running a Reth Node

This is a unified implementation of the Reth node setup that supports running both OP Reth and Base Reth with Flashblocks support.

## Setup

- See hardware requirements mentioned in the master README
- For Base Reth mode: Access to a Flashblocks websocket endpoint (for `RETH_FB_WEBSOCKET_URL`)
  - We provide public websocket endpoints for mainnet and devnet, included in `.env.mainnet` and `.env.sepolia`

## Node Type Selection

Use the `NODE_TYPE` environment variable to select the implementation:

- `NODE_TYPE=vanilla` - OP Reth implementation (default)
- `NODE_TYPE=base` - Base L2 Reth implementation with Flashblocks support

## Environment Variables

| Variable | Purpose | Allowed values | Required | Default | Example |
| --- | --- | --- | --- | --- | --- |
| `NODE_TYPE` | Selects Reth implementation | `vanilla`, `base` | No | `vanilla` (when using the root `docker-compose.yml`) | `NODE_TYPE=base` |
| `CLIENT` | Chooses which client Dockerfile to build via the root `docker-compose.yml` | `reth`, `geth`, `nethermind` | No | `geth` | `CLIENT=reth` |
| `RETH_FB_WEBSOCKET_URL` | Flashblocks websocket endpoint (Base mode only) | valid ws(s) URL | Only when `NODE_TYPE=base` | — | `RETH_FB_WEBSOCKET_URL=wss://...` |

Notes:
- `CLIENT` is consumed by the root `docker-compose.yml` to select the client build context.
- `RETH_FB_WEBSOCKET_URL` is used only when `NODE_TYPE=base`; it is not needed for `vanilla`.

## Running the Node

The node follows the standard `docker-compose` workflow in the master README.

```bash
# Run OP Reth node
CLIENT=reth docker-compose up

# Run Base L2 Reth node with Flashblocks support
NODE_TYPE=base CLIENT=reth docker-compose up
```

## Testing Flashblocks RPC Methods

When running in Base mode (`NODE_TYPE=base`), you can query a pending block using the Flashblocks RPC:

```bash
curl -X POST \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["pending", false],"id":1}' \
  http://localhost:8545
```

## Additional RPC Methods

For a complete list of supported RPC methods, refer to:

- [Standard Ethereum JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Flashblocks RPC Methods](https://docs.base.org/chain/flashblocks#rpc-api) (Base mode only)
