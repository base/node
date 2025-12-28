# Peering configuration guide

This guide provides tips on configuring peer connections for Base nodes.

## Finding peers

- Use `admin_peers` on the execution client to see current peers.
- Share your node's ENR with trusted peers to establish direct connections.

## Adjusting peer limits

- The `maxpeers` flag controls the maximum number of peers your node will accept.
- Lower values reduce bandwidth usage; higher values improve network robustness.

## NAT and firewalls

- Open TCP/UDP ports 30303–30305 on your router for inbound connections.
- Ensure your node advertises the correct public IP by setting `--nat extip:<IP>`.

Keep this document updated as the network evolves.
