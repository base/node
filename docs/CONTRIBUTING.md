# Contributing to the Base node repository

Thank you for considering a contribution!

## Getting started

1. Fork the repository and clone your fork.
2. Run `docker compose up` to start a node locally.
3. Use a new branch per change (`git checkout -b feature/your-description`).

## Code style

- Shell scripts should use `#!/usr/bin/env bash` and `set -euo pipefail`.
- Go code follows `gofmt` style; run `go fmt ./...` before committing.
- Document any new files or APIs.

## Testing your changes

- If you modify the compose setup, run `docker compose config` to ensure the YAML is valid.
- For Go changes, run `go test ./...`.
- For docs, ensure Markdown renders properly and links are correct.

## Opening a pull request

Describe:

- The purpose of your change.
- How you tested it.
- Any related issues it addresses.

Maintainers will review and provide feedback. Welcome aboard!
