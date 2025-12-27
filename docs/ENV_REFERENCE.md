# Environment variable reference

This document lists commonly used environment variables when running a Base node.

## Core variables

- `L1_RPC_URL` – Ethereum L1 RPC endpoint
- `OP_NODE_RPC_URL` – OP node RPC endpoint
- `JWT_SECRET` – JWT secret for engine API authentication

## Optional tuning

- `LOG_LEVEL` – logging verbosity (info, debug)
- `RETH_DB_PATH` – custom path for reth database
- `GOMAXPROCS` – CPU tuning for Go services

Update this document when introducing new configuration options.
