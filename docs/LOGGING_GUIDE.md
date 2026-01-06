# Logging configuration guide

This guide describes how to customise log output for Base node components.

## Log levels

- `error` – show only errors
- `warn` – warnings and errors
- `info` – default verbosity
- `debug` – verbose logs for debugging

Set environment variables to adjust levels:

```bash
LOG_LEVEL=debug
OP_NODE_LOG_LEVEL=info
