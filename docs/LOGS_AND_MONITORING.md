# Logs and monitoring guide

This document provides a short guide for inspecting logs and basic
monitoring signals for a Base node deployed using this repository.

It is intended to complement the Quick Start and Troubleshooting
sections in the README and the official documentation.

---

## 1. Understanding the components

A typical deployment from this repository runs several containers,
including:

- an L2 execution client (for example, `reth`, `geth`, or `nethermind`)
- the OP node process
- supporting services defined in `docker-compose.yml`

When debugging issues, it is often helpful to know which container is
responsible for which part of the stack and to focus on the relevant
logs.

---

## 2. Viewing container logs

From the host machine where Docker is running, you can use
`docker compose` to inspect logs.

Common patterns:

- tail logs from all services:

  ```bash
  docker compose logs -f
