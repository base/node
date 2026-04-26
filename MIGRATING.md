\# Migrating from Geth to Reth



This guide covers migrating an existing Base node from `op-geth` to

`base-reth-node` ahead of the Base V1 upgrade. After V1 activates on

mainnet (early May 2026), only `base-reth-node` and `base-consensus`

will follow the canonical chain — nodes still running `op-geth` will

stop syncing at the activation block.



\## Before you begin



\- Confirm you are running a version of \[base/node](https://github.com/base/node) \*\*before\*\* V1 activates.

\- Ensure you have enough disk space for the Reth data directory (same order of magnitude as your existing Geth data).

\- Read through all steps before starting. The migration requires a brief node downtime.



\## Step 1 — Stop the node



```bash

docker compose down

```



\## Step 2 — Pull the latest base/node



```bash

git pull origin main

```



Confirm the new default client is `reth`:



```bash

grep 'CLIENT:-' docker-compose.yml

\# Expected: dockerfile: ${CLIENT:-reth}/Dockerfile

```



\## Step 3 — Set CLIENT=reth in your env file



Open `.env.mainnet` (or `.env.sepolia` for testnet) and ensure the

following line is present and not commented out:



```bash

CLIENT=reth

```



Alternatively, export it in your shell for the current session:



```bash

export CLIENT=reth

```



\## Step 4 — Bootstrap from a snapshot (recommended)



Syncing Reth from genesis takes days. Use an official snapshot to

bootstrap in hours instead.



Download the latest snapshot:



```bash

\# Mainnet

curl -L https://storage.googleapis.com/base-snapshots/mainnet/latest/reth.tar.gz \\

&#x20; | tar -xz -C /path/to/your/reth-data



\# Sepolia

curl -L https://storage.googleapis.com/base-snapshots/sepolia/latest/reth.tar.gz \\

&#x20; | tar -xz -C /path/to/your/reth-data

```



> \*\*Snapshot freshness\*\*: Official snapshots are updated every 24–48 hours.

> After extraction, Reth will sync the remaining blocks automatically.

> On a reliable connection this catch-up typically takes 1–4 hours.



Set the data directory in your env file:



```bash

RETH\_DATA\_DIR=/path/to/your/reth-data

```



\## Step 5 — Switch to base-consensus



Set `USE\_BASE\_CONSENSUS=true` in your env file. This replaces `op-node`

with `base-consensus`, which is required for V1:



```bash

USE\_BASE\_CONSENSUS=true

```



\## Step 6 — Start the node



```bash

docker compose up -d

```



Follow the logs to confirm both services start cleanly:



```bash

docker compose logs -f

```



\## Step 7 — Verify the migration



Confirm the execution client is now Reth:



```bash

curl -s -X POST http://localhost:8545 \\

&#x20; -H "Content-Type: application/json" \\

&#x20; -d '{"jsonrpc":"2.0","method":"web3\_clientVersion","params":\[],"id":1}' \\

&#x20; | grep -o '"result":"\[^"]\*"'

\# Expected output includes: reth and base

```



Confirm sync is progressing:



```bash

curl -s -X POST http://localhost:8545 \\

&#x20; -H "Content-Type: application/json" \\

&#x20; -d '{"jsonrpc":"2.0","method":"eth\_blockNumber","params":\[],"id":1}'

\# Run twice 30 seconds apart — block number should increase

```



Check peer count (should be > 0 within a few minutes):



```bash

curl -s -X POST http://localhost:8545 \\

&#x20; -H "Content-Type: application/json" \\

&#x20; -d '{"jsonrpc":"2.0","method":"net\_peerCount","params":\[],"id":1}'

```



\## Troubleshooting



\*\*Sync stalled immediately after snapshot extraction\*\*



The snapshot may have been extracted to the wrong directory. Confirm

`RETH\_DATA\_DIR` in your env file points to the directory that contains

the `db/` subdirectory created by the extraction.



\*\*`web3\_clientVersion` still returns `geth`\*\*



`CLIENT=reth` is not being picked up. Verify the variable is set in

your env file and not overridden by a shell export, then run

`docker compose down \&\& docker compose up -d` to force a rebuild.



\*\*Consensus service exits immediately\*\*



Ensure `USE\_BASE\_CONSENSUS=true` is set. Check logs with

`docker compose logs node` for configuration errors, particularly

missing `BASE\_NODE\_\*` variables.



\*\*Peers not connecting after 10+ minutes\*\*



Check that TCP/UDP port 30303 (execution P2P) and 9222 (consensus P2P)

are open in your firewall. If running behind NAT, set `HOST\_IP` to

your public IPv4 address in your env file.



\## Rollback



If you need to revert before V1 activates:



```bash

docker compose down

git checkout <previous-tag>    # e.g. git checkout v0.14.9

export CLIENT=geth

docker compose up -d

```



> After V1 activates, rollback to Geth is not possible — the network

> will have moved past the activation block and Geth will not follow.



\## Environment variable reference



| Variable | Default | Description |

|---|---|---|

| `CLIENT` | `reth` | Execution client (`reth` or `geth`) |

| `RETH\_DATA\_DIR` | `/data` | Path to Reth data directory |

| `USE\_BASE\_CONSENSUS` | `false` | Use `base-consensus` instead of `op-node` |

| `HOST\_IP` | \_(auto)\_ | Public IPv4 for NAT traversal |

| `GETH\_DATA\_DIR` | `/data` | Path to legacy Geth data (not needed after migration) |



\## Additional resources



\- \[Base V1 specification](https://specs.base.org/upgrades/v1/overview)

\- \[Official snapshots](https://docs.base.org/chain/node-operators/snapshots)

\- \[base/node releases](https://github.com/base/node/releases)

\- \[#🛠｜node-operators Discord channel](https://discord.gg/buildonbase)

