# Base node maintenance checklist

This document provides a practical, operator-focused checklist for running and maintaining a Base node.

It is intended as a companion to the official documentation and the main README of this repository, not a replacement for them.

---

## 1. Daily checks

These checks are quick and help you catch issues early.

1. **Sync and health**

   - Confirm that the node is syncing and up to date with the network tip.
   - Use your preferred JSON-RPC client to query:
     - `eth_blockNumber` on the execution client
     - any available sync status endpoint on the op-node (for example, an `optimism_syncStatus`-like method, where supported)
   - Compare the reported block height or timestamps with a trusted public explorer or reference RPC.

2. **Logs**

   - Inspect logs for the execution client and op-node:
     - look for repeated errors
     - check for connection issues to L1 RPCs
     - watch for resource exhaustion warnings (disk, memory, file descriptors)
   - If you use Docker, `docker compose logs` (with appropriate service names) is usually sufficient for a quick view.

3. **Disk usage**

   - Check that the data volume still has enough free space to accommodate chain growth and snapshots.
   - If possible, monitor both:
     - filesystem-level free space
     - size of the actual node data directory.

4. **Basic RPC functionality**

   - Run a simple RPC call (for example, `eth_chainId` and a basic `eth_call` or `eth_getBalance`) to confirm the node is responding correctly.
   - If the node is behind a load balancer or reverse proxy, test through the same path your applications use.

---

## 2. Weekly checks

Once a week, it is useful to do a deeper pass:

1. **Resource usage and trends**

   - Review CPU, memory and network usage over the past week.
   - Confirm that usage patterns match your expectations (for example, no steady upward drift in memory usage).

2. **Snapshots and backups**

   - Verify that any snapshot or backup processes you rely on are completing successfully.
   - Confirm that you can restore from a recent snapshot in a non-production environment.

3. **Configuration drift**

   - Compare your current configuration (`.env`, Docker settings, orchestration configs) against your desired baseline.
   - Make sure ad-hoc changes made during debugging are either reverted or documented.

4. **Software versions**

   - Check for new releases of:
     - this repository (Base node images)
     - your execution client (reth, geth, nethermind)
     - any additional tooling you rely on.
   - Read release notes and decide whether an upgrade is appropriate.

---

## 3. Before performing an upgrade

When planning to upgrade your Base node:

1. **Read release notes carefully**

   - Look for breaking changes, configuration migrations or new environment variables.
   - Pay particular attention to sections describing node operators and infrastructure.

2. **Plan a maintenance window**

   - If this node is part of a production setup or behind a load balancer, schedule a time when it can be taken out of rotation.
   - Ensure that other nodes or fallback infrastructure can handle traffic while this instance is being upgraded.

3. **Take a snapshot or backup**

   - Create a snapshot or backup of the node data and configuration before the upgrade.
   - Verify that the snapshot is complete and accessible.

4. **Prepare rollback steps**

   - Decide in advance how you will roll back if the upgrade exposes unexpected issues:
     - previous container images
     - previous configuration
     - snapshot restore plan.

---

## 4. After an upgrade

Once the node has been upgraded:

1. **Verify sync and health**

   - Confirm that the node starts successfully and resumes syncing.
   - Re-run your health and sync checks (block height, RPC probes, logs).

2. **Validate configuration**

   - Make sure new configuration options are set correctly.
   - Confirm that deprecated options are either removed or updated.

3. **Run a small functional test**

   - If possible, perform a small end-to-end test using your normal application workflow:
     - connect via RPC
     - send a simple transaction
     - confirm it is processed as expected.

4. **Monitor closely**

   - Monitor logs and metrics more closely than usual during the first hours after the upgrade.
   - Watch for regressions in performance, memory usage, or error rates.

---

## 5. Incident response quick reference

If your Base node appears unhealthy or out of sync:

1. **Take it out of rotation**

   - If the node is serving production traffic behind a load balancer, remove it from the pool to avoid impacting users.

2. **Collect information**

   - Capture:
     - recent logs
     - resource utilisation snapshots
     - configuration files or environment (with secrets removed)
     - observations from your monitoring system.

3. **Check for known issues**

   - Review:
     - recent releases
     - open issues or discussions related to the symptoms you are seeing.

4. **Decide between repair and rebuild**

   - For some cases, it may be faster and safer to:
     - stop the node
     - restore from a known-good snapshot
     - or resync from scratch,
     rather than attempting to repair a potentially corrupted state.

5. **Document and share learnings**

   - Keep a short incident log so that future maintenance and automation can prevent similar issues.
   - When appropriate, share anonymized details with maintainers as part of a bug report.

---

## 6. Notes

This checklist is intentionally generic and should be adapted to your environment:

- the exact commands you use may differ depending on how you deploy the node
- you may have additional compliance, security, or observability requirements to follow.

Treat this document as a starting point and refine it as you gain operational experience with your Base nodes.
