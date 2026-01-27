# Contributing Guidelines

Thank you for your interest in contributing to Base Node! This document
provides best practices and instructions for making contributions.

## Getting started

1. Fork the repository and clone your fork.
2. Create a new branch for your change.
3. Install dependencies as described in the README.
4. Ensure you have Go 1.20 or later installed.

## Coding standards

- Follow the existing project structure and import paths.
- Run `go vet` and `golangci-lint` before submitting a PR.
- Write unit tests for new functions or significant changes.

## Submitting changes

1. Keep your PRs focused and small.
2. Include a clear description of what your change does.
3. Update documentation if your change affects behaviour.
4. Ensure all tests pass: `go test ./...`.

We appreciate your contributions and reviews help keep the codebase healthy!
