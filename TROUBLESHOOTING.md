
# Base Node Troubleshooting Guide

This guide covers common issues encountered when operating a Base node (using `op-geth`, `reth`, or `nethermind` paired with `op-node`).

## Table of Contents
1. [Syncing & Performance](#1-syncing--performance)
   - [Context Deadline Exceeded / Execution Failed](#context-deadline-exceeded--payload-execution-failed)
   - [Node Not Syncing / Stuck at Block 0](#node-not-syncing--stuck-at-block-0)
   - [Snap Sync Stuck "Healing"](#snap-sync-stuck-healing)
2. [Configuration & Connectivity](#2-configuration--connectivity)
   - [Engine API / JWT Token Errors](#engine-api--jwt-token-errors-401-unauthorized)
   - [L1 RPC Connection Failures](#l1-rpc-connection-failures)
   - [Peer Discovery Issues (0 Peers)](#peer-discovery-issues)
3. [Docker & Permissions](#3-docker--permissions)
   - [Permission Denied on Volumes](#permission-denied-on-volumes)

---

## 1. Syncing & Performance

### Context Deadline Exceeded / Payload Execution Failed

**Error Log:**
```text
Payload execution failed: context deadline exceeded

```

**Diagnosis:**
This is the most common error for node operators. While it looks like a network timeout, it is almost a **High Assurance indicator of a Disk I/O bottleneck**. The execution client cannot read/write the state trie fast enough to keep up with the chain tip.

**Solution:**

1. **Hardware Check:** Ensure you are running on an **NVMe SSD**. SATA SSDs and HDDs are insufficient for Base Mainnet.
2. **Cloud Storage:** If using AWS/GCP, standard volumes (e.g., AWS gp2) will fail. You must use **Provisioned IOPS** (io1/io2) or high-throughput tiers (gp3 with maxed throughput).
3. **IOPS Requirement:** Ensure your disk provides **10,000+ IOPS**.

### Node Not Syncing / Stuck at Block 0

**Symptoms:**

* `op-node` logs show `L2 output root mismatch` or genesis hash mismatch.
* Node refuses to advance past the genesis block.

**Solution:**

1. **Check Network:** Ensure you are using the correct `--network` flag (e.g., `mainnet` vs `sepolia`).
2. **Genesis File:** If manually configuring, verify your `genesis.json` matches the official Base genesis for your targeted network.
3. **Rollup Config:** Ensure `rollup.json` is correctly loaded and matches the network constraints.

### Snap Sync Stuck "Healing"

**Symptoms:**

* Geth/Nethermind logs show "Healing state" for an extended period (hours/days) without completing.

**Solution:**

1. **Patience:** "Healing" involves downloading the state trie leaves. On a chain as large as Base, this can take significant time depending on network speed.
2. **Restart with Fresh DB:** If the healing phase hangs for >24h or the database is corrupted, it is often faster to wipe the `data` directory and restart the snap sync from scratch.
3. **Peers:** Ensure you have enough healthy peers (20+) to source state data from.

---

## 2. Configuration & Connectivity

### Engine API / JWT Token Errors (401 Unauthorized)

**Error Log:**

```text
Engine API request failed: 401 Unauthorized

```

or

```text
Failed to unmarshal JWT secret

```

**Diagnosis:**
The `op-node` (consensus client) cannot authenticate with the execution client (geth/reth) because the JWT secrets do not match.

**Solution:**

1. **Shared Secret:** Ensure both containers are pointing to the **exact same** `jwt.hex` file.
2. **File Path:** In Docker, ensure the volume mount for the JWT file is correct in `docker-compose.yml`. Both services must be able to read the file from their respective mounted paths.
3. **Generate New Secret:** If unsure, regenerate the token:
```bash
openssl rand -hex 32 > jwt.hex

```



### L1 RPC Connection Failures

**Error Log:**

```text
Failed to fetch L1 block info
connection refused

```

**Diagnosis:**
The `op-node` requires a reliable connection to an Ethereum L1 node (RPC) to derive the L2 chain.

**Solution:**

1. **Check L1 Endpoint:** Ensure `OP_NODE_L1_ETH_RPC` points to a valid, synced Ethereum Mainnet (or Sepolia) node.
2. **Rate Limits:** If using a public RPC (like Infura/Alchemy free tier), you may be hitting rate limits. Base derivation requires frequent polling. Use a paid plan or run your own L1 node.
3. **Consensus Client:** Ensure your L1 RPC supports the Engine API if you are running a full L1 node stack.

### Peer Discovery Issues

**Symptoms:**

* `peer count=0`
* Logs showing "Discovered 0 nodes"

**Solution:**

1. **Port Forwarding:** Ensure ports `30303` (TCP/UDP) for execution and `9003` (TCP/UDP) for consensus are open and forwarded on your firewall/router.
2. **Bootnodes:** Check if your client is connecting to the hardcoded bootnodes. If not, manually add Base bootnodes to your configuration.
3. **Static Peers:** If discovery fails, add trusted static peers to your config to jumpstart the connection.

---

## 3. Docker & Permissions

### Permission Denied on Volumes

**Error Log:**

```text
open /datadir/geth/chaindata/LOCK: permission denied
```

**Diagnosis:**
The Docker container user (usually non-root) does not have write access to the mounted host directory.

**Solution:**

1. **Chown Directory:** On the host machine, change ownership of the data directory to the user ID expected by the container (often `1000:1000` or `1001:1001`).
```bash
sudo chown -R 1000:1000 ./data-directory
```


2. **User Flag:** explicitly set the `user: "1000:1000"` in your `docker-compose.yml` if your host user matches that ID.
