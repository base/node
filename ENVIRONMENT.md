# Environment configuration for Base Node

This document provides a brief overview of environment variables used when running a Base node with Docker and `docker compose`. For a full guide, please refer to the official Base documentation: https://docs.base.org/base-chain/node-operators/run-a-base-node

## L1 connectivity

These variables configure how `op-node` connects to Ethereum L1:

- `OP_NODE_L1_ETH_RPC` – HTTP(S) endpoint of your Ethereum L1 full node or RPC provider.
- `OP_NODE_L1_BEACON` – Beacon chain endpoint for your Ethereum L1 node.
- `OP_NODE_L1_RPC_KIND` – Describes the type of L1 provider (for example `basic` or `erigon`), matching the values documented in the OP Stack docs.

When using `docker-compose`, these values are typically set in `.env.mainnet` or `.env.sepolia` depending on which network you are targeting.

## L2 / execution clients

Depending on which execution client you run, different variables are relevant:

- `RETH_CHAIN` – Chain identifier for Reth when used as the execution client.
- `RETH_SEQUENCER_HTTP` – URL of the sequencer RPC endpoint that Reth should use.
- `OP_GETH_*` and `NETHERMIND_*` – Variables that control which op-geth or Nethermind versions are used in the Docker images, as pinned in `versions.env`.

## Recommended references

For more details on how to configure and tune your node, see:

- Base node getting started guide: https://docs.base.org/base-chain/node-operators/run-a-base-node
- Node performance guide: https://docs.base.org/base-chain/node-operators/performance-tuning
- Troubleshooting documentation: https://docs.base.org/base-chain/node-operators/troubleshooting
