# Base node monitoring guide

This document provides a practical overview of what to monitor when running a Base node in production or pre-production environments.

It is meant to complement the main README and external documentation, not replace them.

---

## 1. High-level goals

Monitoring for a Base node should help you answer at least the following questions:

- Is the node healthy and reachable over RPC?
- Is the node reasonably close to the network tip?
- Are resources (CPU, memory, disk, network) within expected ranges?
- Are there errors in the logs that require attention?
- Are recent upgrades or configuration changes behaving as expected?

---

## 2. Core health checks

### 2.1 RPC availability

From a monitoring job or script, periodically:

- call `eth_blockNumber` on the execution client RPC
- optionally, call any available sync status method on the op-node (for example, an `optimism_syncStatus`-style endpoint, where supported)

Basic checks:

- the RPC endpoint responds within an acceptable latency
- the response is well-formed JSON
- the reported block number is incrementing over time

### 2.2 Sync distance to tip

Compare the local L2 block height or timestamp with:

- a trusted public Base RPC
- or an explorer that exposes the current network tip

Alert if the gap exceeds a threshold that is acceptable for your use case (for example, a few dozen blocks or a small number of minutes).

---

## 3. Resource monitoring

### 3.1 CPU and memory

Track:

- average and peak CPU usage for:
  - the execution client container (reth / geth / nethermind)
  - the op-node container
- memory usage (RSS) for the same processes or containers

Typical guidelines:

- sustained 100% CPU usage for long periods may indicate under-provisioned hardware or misconfiguration
- steady, unbounded growth in memory usage should be investigated

### 3.2 Disk and I/O

Monitor:

- free space on the filesystem where node data is stored
- growth rate of the data directory
- disk I/O saturation (queue depth, read/write latency) if available

Set alerts for low disk space well before the node runs out (for example, below 20–25% free space, depending on your risk tolerance).

---

## 4. Logs

### 4.1 Error and warning patterns

Configure log collection (for example, via Docker logging drivers, `journald`, or an external collector) and watch for:

- repeated connection errors to:
  - L1 RPC
  - beacon endpoints
  - sequencer endpoints
- errors related to:
  - consensus
  - block import
  - state corruption
- restarts or crashes of the node processes

### 4.2 Log levels and rotation

Good practices:

- in steady-state production, avoid excessively verbose log levels that can:
  - increase disk usage
  - make it harder to spot important messages
- ensure log rotation is configured:
  - limit the size and number of log files
  - avoid filling the disk with historical logs

---

## 5. Configuration and environment drift

When operating more than one node:

- regularly compare configuration files (`.env`, Docker settings, orchestration manifests) across instances
- ensure that:
  - nodes that are supposed to be identical actually run with the same settings
  - experimental flags or temporary changes are documented and eventually cleaned up

Consider periodically exporting and storing:

- hashes of configuration files
- container image versions
- selected environment variables (without secrets)

---

## 6. Alerts and dashboards

A minimal set of alerts for a Base node might include:

- **RPC availability:** endpoint unreachable or error rate above a threshold
- **Sync lag:** node is behind the tip by more than a configured number of blocks or minutes
- **Disk space:** free space below a configured threshold on the data volume
- **Node restarts:** abnormal restart rate over a given time window
- **Error spikes:** sudden increase in error-level log entries

Dashboards can visualize:

- block height over time
- resource usage
- key latency metrics (RPC, database, network)

---

## 7. After upgrades and configuration changes

After upgrading the node or changing important configuration:

1. Monitor more closely for a period of time:
   - resource usage
   - sync lag
   - error logs
2. Compare metrics before and after the change where possible.
3. Be ready to roll back to a previous configuration or snapshot if new issues appear.

---

## 8. Adapting this guide

Every environment is different. Use this document as a starting point and adapt it to:

- your hardware and scaling strategy
- your observability stack (Prometheus, hosted monitoring, etc.)
- your operational requirements and risk appetite

Additional guidance specific to Base node configuration can be found in the main README and the wider Base documentation.
