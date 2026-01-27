# Local Logging & Troubleshooting Guide

This document describes a few small helpers for working with a local Base node.

## 1. Check running services

To see which services are up, run:

    docker compose ps

This helps you confirm that the node and execution clients are running.

## 2. Tail node logs

Use the helper script to follow logs:

    ./scripts/tail_node_logs.sh

Or tail a specific service only (for example, the op-node):

    ./scripts/tail_node_logs.sh op-node

## 3. Quick RPC health check

You can verify that the JSON-RPC endpoint is responding with a simple script:

    ./scripts/check_node_health.sh
    ./scripts/check_node_health.sh http://localhost:8545

If the script prints a latest block number, your node is responding to `eth_blockNumber`.
If it reports an HTTP error or empty result, check your environment variables and Docker logs.
