# Local snapshot restore guide

This document describes a practical workflow for restoring a Base node
from a snapshot using the configuration provided in this repository.

It is intended as a complement to the official documentation and
troubleshooting guides.

---

## 1. Prerequisites

Before restoring from a snapshot, make sure that:

- you have cloned the `base/node` repository and are working from its
  root directory
- you are using a supported Docker and Docker Compose installation
- you have configured your L1 RPC endpoints in one of:
  - `.env.mainnet` for mainnet
  - `.env.sepolia` for testnet
- you have enough free disk space for:
  - the current chain data,
  - the snapshot archive,
  - and additional headroom for future growth

For the latest hardware and storage recommendations, refer to the
official Base documentation.

---

## 2. Choosing the network and environment file

Decide which network you want to restore:

- **Mainnet**:
  - use `.env.mainnet`
- **Sepolia testnet**:
  - use `.env.sepolia`

Verify that the appropriate file contains valid values for at least:

- `OP_NODE_L1_ETH_RPC`
- `OP_NODE_L1_BEACON`
- `OP_NODE_L1_BEACON_ARCHIVER`

These must point to reachable and fully synced Ethereum L1 endpoints.

---

## 3. Stopping the running node

If you already have a node running from this repository, stop it before
restoring data.

From the repository root:

1. List running containers related to the Base node (for example using
   `docker compose ps`).
2. Stop the stack:

    - for default configuration:

      - `docker compose down`

    - or, if you are using environment overrides:

      - `NETWORK_ENV=.env.sepolia docker compose down`
      - `CLIENT=reth docker compose down`

Make sure all containers from this stack have stopped before proceeding.

---

## 4. Downloading the snapshot

Snapshots are published by the Base team and referenced from the
official documentation.

Typical steps:

1. Visit the snapshot information page referenced in the README or docs.
2. Choose a snapshot that matches:
   - your target network (mainnet or sepolia)
   - the client you plan to run (for example `reth` or `geth`)
3. Download the archive to a directory on the same filesystem where your
   node data will live.

For example, you might end up with a file such as:

- `/data/base/snapshots/reth-mainnet-2025-01-01.tar.zst`

The exact filename and compression format depend on how snapshots are
published at the time you are setting up your node.

---

## 5. Preparing the data directory

Identify the data directory used by your chosen client. This is
typically configured in:

- `docker-compose.yml`
- and/or environment variables in `.env.mainnet` or `.env.sepolia`

Common patterns include paths under:

- `./data/reth`
- `./data/geth`
- `./data/nethermind`

Before restoring the snapshot:

1. Stop all containers (as described in section 3).
2. Make a backup of any existing data directory if you want to keep it.
3. Remove or empty the data directory where the snapshot will be
   restored.

---

## 6. Restoring the snapshot

Use the decompression tool that matches the snapshot format (for
example, `tar`, `zstd`, or similar).

A generic pattern looks like:

- `cd /path/to/base-node-repo`
- `cd /path/to/data-directory-for-your-client`
- extract the snapshot archive into this directory

Example (adjust for your actual path and tools):

- `tar --extract --file /data/base/snapshots/reth-mainnet-2025-01-01.tar.zst`
- or, if you need a decompression step, run the appropriate `zstd` or
  `unzstd` command first and then `tar`.

After extraction, the data directory should contain the database files
expected by your chosen client.

---

## 7. Starting the node with restored data

Once the snapshot data is in place:

1. Return to the repository root.
2. Start the node with the appropriate environment:

   - for mainnet with default client:

     - `docker compose up --build`

   - for testnet:

     - `NETWORK_ENV=.env.sepolia docker compose up --build`

   - for a specific client:

     - `CLIENT=reth docker compose up --build`

3. Monitor the logs to confirm that:
   - the client starts successfully,
   - the node proceeds from the restored snapshot height and continues
     syncing.

---

## 8. Verifying the restored node

To verify that the node is healthy after a snapshot restore:

- check that the latest block number is moving forward
- confirm that:
  - JSON-RPC endpoints respond as expected
  - peer counts and sync status are in a reasonable range
- compare the reported chain head with a trusted public RPC endpoint

If you see repeated errors in the logs, or the node does not progress,
double-check:

- the integrity of the snapshot archive
- that the snapshot matches:
  - the correct network (mainnet vs testnet)
  - the client implementation you are running

---

## 9. Troubleshooting

If you encounter issues while restoring from a snapshot:

- review the node logs for your client containers
- verify that:
  - L1 RPC endpoints are reachable and correctly configured
  - disk space and I/O performance are sufficient
- consult the official troubleshooting guide linked from the README

If the problem persists, consider opening a GitHub issue or asking for
help in the Base node operators channel, including:

- details of the snapshot used
- your client type and version
- relevant excerpts from your logs
