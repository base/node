# RPC Troubleshooting Notes

This document collects a few quick checks to perform when your Base node
is running but RPC calls do not behave as expected.

## 1. Check that ports are exposed

If you are running via Docker, ensure that the HTTP port is mapped
from the container to the host, for example:

- `8545:8545` for the execution client
- `8547:8547` for the OP node (if applicable)

You can verify open ports via:

    docker ps
    docker inspect <container>

## 2. Sanity-check JSON-RPC

Use a simple `curl` request:

    curl -X POST http://localhost:8545 \
      -H "Content-Type: application/json" \
      -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

If the node is healthy, you should get a hex-encoded block number.

## 3. Confirm environment variables

Common variables that affect connectivity:

- `L1_RPC_URL`
- `OP_NODE_RPC_URL`
- `CLIENT`
- `HOST_DATA_DIR`

If you change any of these, rebuild or restart your containers to apply
the new configuration.

## 4. Logs

Always check the logs of both the execution client and the OP node.
Look for:

- repeated connection failures,
- invalid RPC URLs,
- authentication or rate limit errors.

They will usually point you directly to the misconfiguration that needs
to be fixed.
