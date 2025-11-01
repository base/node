# ===============================
# Stage 1: Build op-node
# ===============================
FROM golang:1.24 AS op
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

RUN curl -sSfL 'https://just.systems/install.sh' | bash -s -- --to /usr/local/bin

WORKDIR /app
COPY versions.env /tmp/versions.env

RUN set -e && \
    source /tmp/versions.env && \
    test -n "$OP_NODE_REPO" && \
    test -n "$OP_NODE_TAG" && \
    git clone "$OP_NODE_REPO" --branch "$OP_NODE_TAG" --single-branch . && \
    git switch -c "branch-$OP_NODE_TAG" && \
    bash -c '[ "$(git rev-parse HEAD)" = "$OP_NODE_COMMIT" ]' && \
    make -C op-node VERSION="$OP_NODE_TAG" op-node

# ===============================
# Stage 2: Build geth
# ===============================
FROM golang:1.24 AS geth
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /app
COPY versions.env /tmp/versions.env

RUN set -e && \
    source /tmp/versions.env && \
    test -n "$OP_GETH_REPO" && \
    test -n "$OP_GETH_TAG" && \
    git clone "$OP_GETH_REPO" --branch "$OP_GETH_TAG" --single-branch . && \
    git switch -c "branch-$OP_GETH_TAG" && \
    bash -c '[ "$(git rev-parse HEAD)" = "$OP_GETH_COMMIT" ]' && \
    go run build/ci.go install -static ./cmd/geth

# ===============================
# Stage 3: Final image
# ===============================
FROM ubuntu:24.04
LABEL maintainer="Your Name <you@example.com>" \
      description="Combined op-node + geth environment with supervisord"

RUN apt-get update && \
    apt-get install -y --no-install-recommends jq curl supervisor ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /var/log/supervisor /app

WORKDIR /app

# Copy built binaries
COPY --from=op /app/op-node/bin/op-node ./op-node
COPY --from=geth /app/build/bin/geth ./geth

# Copy configuration
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY geth/geth-entrypoint ./execution-entrypoint
COPY op-node-entrypoint .

RUN chmod +x ./execution-entrypoint ./op-node-entrypoint ./op-node ./geth

EXPOSE 8545 8551 6060

ENTRYPOINT ["/usr/bin/supervisord"]
