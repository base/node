name: Pull Request

on:
  pull_request:
  workflow_dispatch:

permissions:
  contents: read

jobs:
  geth:
    strategy:
      matrix:
        settings:
          - arch: linux/amd64
            runs-on: ubuntu-24.04
          - arch: linux/arm64
            runs-on: ubuntu-24.04-arm
    runs-on: ${{ matrix.settings.runs-on }}
    steps:
      - name: Harden the runner (Audit all outbound calls)
        uses: step-security/harden-runner@002fdce3c6a235733a90a27c80493a3241e56863 # v2.12.1
        with:
          egress-policy: audit

      - name: Checkout
        uses: actions/checkout@f43a0e5ff2bd294095638e18286ca9a3d1956744 # v3.6.0
        with:
          ref: ${{ github.event.pull_request.head.sha }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

      - name: Build the Docker image
        uses: docker/build-push-action@263435318d21b8e681c14492fe198d362a7d2c83 # v6.18.0
        with:
          context: .
          file: geth/Dockerfile
          push: false
          platforms: ${{ matrix.settings.arch }}

  reth:
    strategy:
      matrix:
        settings:
          - arch: linux/amd64
            runs-on: ubuntu-24.04
            features: jemalloc,asm-keccak,optimism
          - arch: linux/arm64
            runs-on: ubuntu-24.04-arm
            features: jemalloc,optimism
    runs-on: ${{ matrix.settings.runs-on}}
    steps:
      - name: Harden the runner (Audit all outbound calls)
        uses: step-security/harden-runner@002fdce3c6a235733a90a27c80493a3241e56863 # v2.12.1
        with:
          egress-policy: audit

      - name: Checkout
        uses: actions/checkout@f43a0e5ff2bd294095638e18286ca9a3d1956744 # v3.6.0
        with:
          ref: ${{ github.event.pull_request.head.sha }}
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1
      - name: Build the Docker image
        uses: docker/build-push-action@263435318d21b8e681c14492fe198d362a7d2c83 # v6.18.0
        with:
          context: .
          file: reth/Dockerfile
          push: false
          build-args: |
            FEATURES=${{ matrix.settings.features }}
          platforms: ${{ matrix.settings.arch }}

  nethermind:
    strategy:
      matrix:
        settings:
          - arch: linux/amd64
            runs-on: ubuntu-24.04
          - arch: linux/arm64
            runs-on: ubuntu-24.04-arm
    runs-on: ${{ matrix.settings.runs-on}}
    steps:
      - name: Harden the runner (Audit all outbound calls)
        uses: step-security/harden-runner@002fdce3c6a235733a90a27c80493a3241e56863 # v2.12.1
        with:
          egress-policy: audit

      - name: Checkout
        uses: actions/checkout@f43a0e5ff2bd294095638e18286ca9a3d1956744 # v3.6.0
        with:
          ref: ${{ github.event.pull_request.head.sha }}
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1
      - name: Build the Docker image
        uses: docker/build-push-action@263435318d21b8e681c14492fe198d362a7d2c83 # v6.18.0
        with:
          context: .
          file: nethermind/Dockerfile
          push: false
          platforms: ${{ matrix.settings.arch }}
