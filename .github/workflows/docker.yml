name: Tag Docker image

on:
  push:
    branches:
      - "main"
    tags:
      - "v*"

env:
  REGISTRY: ghcr.io
  NAMESPACE: ghcr.io/base
  GETH_DEPRECATED_IMAGE_NAME: node
  GETH_IMAGE_NAME: node-geth
  RETH_IMAGE_NAME: node-reth
  NETHERMIND_IMAGE_NAME: node-nethermind

permissions:
  contents: read
  packages: write

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

      - name: Log into the Container registry
        uses: docker/login-action@74a5d142397b4f367a81961eba4e8cd7edddf772 # v3.4.0
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata for the Docker image
        id: meta
        uses: docker/metadata-action@818d4b7b91585d195f67373fd9cb0332e31a7175 # v4.6.0
        with:
          images: |
            ${{ env.NAMESPACE }}/${{ env.GETH_DEPRECATED_IMAGE_NAME }}
            ${{ env.NAMESPACE }}/${{ env.GETH_IMAGE_NAME }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

      - name: Build and push the Docker image
        id: build
        uses: docker/build-push-action@263435318d21b8e681c14492fe198d362a7d2c83 # v6.18.0
        with:
          context: .
          file: geth/Dockerfile
          tags: ${{ env.NAMESPACE }}/${{ env.GETH_DEPRECATED_IMAGE_NAME }},${{ env.NAMESPACE }}/${{ env.GETH_IMAGE_NAME }}
          labels: ${{ steps.meta.outputs.labels }}
          platforms: ${{ matrix.settings.arch }}
          outputs: type=image,push-by-digest=true,name-canonical=true,push=true

      - name: Export digest
        run: |
          mkdir -p ${{ runner.temp }}/digests
          digest="${{ steps.build.outputs.digest }}"
          touch "${{ runner.temp }}/digests/${digest#sha256:}"

      - name: Prepare
        run: |
          platform=${{ matrix.settings.arch }}
          echo "PLATFORM_PAIR=${platform//\//-}" >> $GITHUB_ENV
  
      - name: Upload digest
        uses: actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02 # v4.6.2
        with:
          name: digests-geth-${{ env.PLATFORM_PAIR }}
          path: ${{ runner.temp }}/digests/*
          if-no-files-found: error
          retention-days: 1
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
    runs-on: ${{ matrix.settings.runs-on }}
    steps:
      - name: Harden the runner (Audit all outbound calls)
        uses: step-security/harden-runner@002fdce3c6a235733a90a27c80493a3241e56863 # v2.12.1
        with:
          egress-policy: audit

      - name: Checkout
        uses: actions/checkout@ee0669bd1cc54295c223e0bb666b733df41de1c5 # v2.7.0

      - name: Log into the Container registry
        uses: docker/login-action@74a5d142397b4f367a81961eba4e8cd7edddf772 # v3.4.0
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata for the Docker image
        id: meta
        uses: docker/metadata-action@818d4b7b91585d195f67373fd9cb0332e31a7175 # v4.6.0
        with:
          images: |
            ${{ env.NAMESPACE }}/${{ env.RETH_IMAGE_NAME }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

      - name: Build and push the Docker image
        id: build
        uses: docker/build-push-action@263435318d21b8e681c14492fe198d362a7d2c83 # v6.18.0
        with:
          context: .
          file: reth/Dockerfile
          tags: ${{ env.NAMESPACE }}/${{ env.RETH_IMAGE_NAME }}
          labels: ${{ steps.meta.outputs.labels }}
          build-args: |
            FEATURES=${{ matrix.settings.features }}
          platforms: ${{ matrix.settings.arch }}
          outputs: type=image,push-by-digest=true,name-canonical=true,push=true

      - name: Export digest
        run: |
          mkdir -p ${{ runner.temp }}/digests
          digest="${{ steps.build.outputs.digest }}"
          touch "${{ runner.temp }}/digests/${digest#sha256:}"

      - name: Prepare
        run: |
          platform=${{ matrix.settings.arch }}
          echo "PLATFORM_PAIR=${platform//\//-}" >> $GITHUB_ENV

      - name: Upload digest
        uses: actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02 # v4.6.2
        with:
          name: digests-reth-${{ env.PLATFORM_PAIR }}
          path: ${{ runner.temp }}/digests/*
          if-no-files-found: error
          retention-days: 1

  nethermind:
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
        uses: actions/checkout@ee0669bd1cc54295c223e0bb666b733df41de1c5 # v2.7.0

      - name: Log into the Container registry
        uses: docker/login-action@74a5d142397b4f367a81961eba4e8cd7edddf772 # v3.4.0
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

      - name: Extract metadata for the Docker image
        id: meta
        uses: docker/metadata-action@818d4b7b91585d195f67373fd9cb0332e31a7175 # v4.6.0
        with:
          images: |
            ${{ env.NAMESPACE }}/${{ env.NETHERMIND_IMAGE_NAME }}

      - name: Build and push the Docker image
        id: build
        uses: docker/build-push-action@263435318d21b8e681c14492fe198d362a7d2c83 # v6.18.0
        with:
          context: .
          file: nethermind/Dockerfile
          tags: ${{ env.NAMESPACE }}/${{ env.NETHERMIND_IMAGE_NAME }}
          labels: ${{ steps.meta.outputs.labels }}
          platforms: ${{ matrix.settings.arch }}
          outputs: type=image,push-by-digest=true,name-canonical=true,push=true

      - name: Export digest
        run: |
          mkdir -p ${{ runner.temp }}/digests
          digest="${{ steps.build.outputs.digest }}"
          touch "${{ runner.temp }}/digests/${digest#sha256:}"

      - name: Prepare
        run: |
          platform=${{ matrix.settings.arch }}
          echo "PLATFORM_PAIR=${platform//\//-}" >> $GITHUB_ENV

      - name: Upload digest
        uses: actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02 # v4.6.2
        with:
          name: digests-nethermind-${{ env.PLATFORM_PAIR }}
          path: ${{ runner.temp }}/digests/*
          if-no-files-found: error
          retention-days: 1


  merge-geth:
    runs-on: ubuntu-latest
    needs:
      - geth
    steps:
      - name: Harden the runner (Audit all outbound calls)
        uses: step-security/harden-runner@002fdce3c6a235733a90a27c80493a3241e56863 # v2.12.1
        with:
          egress-policy: audit

      - name: Download digests
        uses: actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093 # v4.3.0
        with:
          path: ${{ runner.temp }}/digests
          pattern: digests-geth-*
          merge-multiple: true

      - name: Log into the Container registry
        uses: docker/login-action@74a5d142397b4f367a81961eba4e8cd7edddf772 # v3.4.0
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

      - name: Extract metadata for the Docker image
        id: meta
        uses: docker/metadata-action@902fa8ec7d6ecbf8d84d538b9b233a880e428804 # v5.7.0
        with:
          images: |
            ${{ env.NAMESPACE }}/${{ env.GETH_DEPRECATED_IMAGE_NAME }}
            ${{ env.NAMESPACE }}/${{ env.GETH_IMAGE_NAME }}

      - name: Create manifest list and push
        working-directory: ${{ runner.temp }}/digests
        run: |
          docker buildx imagetools create $(jq -cr '.tags | map("-t " + .) | join(" ")' <<< "$DOCKER_METADATA_OUTPUT_JSON") \
            $(printf '${{ env.NAMESPACE }}/${{ env.GETH_DEPRECATED_IMAGE_NAME }}@sha256:%s ' *)
          docker buildx imagetools create $(jq -cr '.tags | map("-t " + .) | join(" ")' <<< "$DOCKER_METADATA_OUTPUT_JSON") \
            $(printf '${{ env.NAMESPACE }}/${{ env.GETH_IMAGE_NAME }}@sha256:%s ' *)

      - name: Inspect image
        run: |
          docker buildx imagetools inspect ${{ env.NAMESPACE }}/${{ env.GETH_DEPRECATED_IMAGE_NAME }}:${{ steps.meta.outputs.version }}
          docker buildx imagetools inspect ${{ env.NAMESPACE }}/${{ env.GETH_IMAGE_NAME }}:${{ steps.meta.outputs.version }}

  merge-reth:
    runs-on: ubuntu-latest
    needs:
      - reth
    steps:
      - name: Harden the runner (Audit all outbound calls)
        uses: step-security/harden-runner@002fdce3c6a235733a90a27c80493a3241e56863 # v2.12.1
        with:
          egress-policy: audit

      - name: Download digests
        uses: actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093 # v4.3.0
        with:
          path: ${{ runner.temp }}/digests
          pattern: digests-reth-*
          merge-multiple: true

      - name: Log into the Container registry
        uses: docker/login-action@74a5d142397b4f367a81961eba4e8cd7edddf772 # v3.4.0
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

      - name: Extract metadata for the Docker image
        id: meta
        uses: docker/metadata-action@902fa8ec7d6ecbf8d84d538b9b233a880e428804 # v5.7.0
        with:
          images: |
            ${{ env.NAMESPACE }}/${{ env.RETH_IMAGE_NAME }}

      - name: Create manifest list and push
        working-directory: ${{ runner.temp }}/digests
        run: |
          docker buildx imagetools create $(jq -cr '.tags | map("-t " + .) | join(" ")' <<< "$DOCKER_METADATA_OUTPUT_JSON") \
            $(printf '${{ env.NAMESPACE }}/${{ env.RETH_IMAGE_NAME }}@sha256:%s ' *)

      - name: Inspect image
        run: |
          docker buildx imagetools inspect ${{ env.NAMESPACE }}/${{ env.RETH_IMAGE_NAME }}:${{ steps.meta.outputs.version }}

  merge-nethermind:
    runs-on: ubuntu-latest
    needs:
      - nethermind
    steps:
      - name: Harden the runner (Audit all outbound calls)
        uses: step-security/harden-runner@002fdce3c6a235733a90a27c80493a3241e56863 # v2.12.1
        with:
          egress-policy: audit

      - name: Download digests
        uses: actions/download-artifact@d3f86a106a0bac45b974a628896c90dbdf5c8093 # v4.3.0
        with:
          path: ${{ runner.temp }}/digests
          pattern: digests-nethermind-*
          merge-multiple: true

      - name: Log into the Container registry
        uses: docker/login-action@74a5d142397b4f367a81961eba4e8cd7edddf772 # v3.4.0
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@e468171a9de216ec08956ac3ada2f0791b6bd435 # v3.11.1

      - name: Extract metadata for the Docker image
        id: meta
        uses: docker/metadata-action@902fa8ec7d6ecbf8d84d538b9b233a880e428804 # v5.7.0
        with:
          images: |
            ${{ env.NAMESPACE }}/${{ env.NETHERMIND_IMAGE_NAME }}

      - name: Create manifest list and push
        working-directory: ${{ runner.temp }}/digests
        run: |
          docker buildx imagetools create $(jq -cr '.tags | map("-t " + .) | join(" ")' <<< "$DOCKER_METADATA_OUTPUT_JSON") \
            $(printf '${{ env.NAMESPACE }}/${{ env.NETHERMIND_IMAGE_NAME }}@sha256:%s ' *)

      - name: Inspect image
        run: |
          docker buildx imagetools inspect ${{ env.NAMESPACE }}/${{ env.NETHERMIND_IMAGE_NAME }}:${{ steps.meta.outputs.version }}
