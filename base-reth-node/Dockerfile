FROM golang:1.22 AS op

WORKDIR /app

ENV REPO=https://github.com/ethereum-optimism/optimism.git
ENV VERSION=v1.13.0
ENV COMMIT=5f003211aed7469eed7df666291a62c025d1c46c
RUN git clone $REPO --branch op-node/$VERSION --single-branch . && \
    git switch -c branch-$VERSION && \
    bash -c '[ "$(git rev-parse HEAD)" = "$COMMIT" ]' || (echo "Commit hash verification failed" && exit 1)

# Install 'just'
RUN curl -sSfL 'https://just.systems/install.sh' | bash -s -- --to /usr/local/bin

RUN cd op-node && \
    make VERSION=$VERSION op-node

FROM rust:1.85 AS reth

WORKDIR /app

RUN apt-get update && apt-get -y upgrade && \
    apt-get install -y git libclang-dev pkg-config curl build-essential && \
    rm -rf /var/lib/apt/lists/*

# Base node-reth repository configuration
ENV REPO=https://github.com/base/node-reth.git
ENV VERSION=main
ENV COMMIT=c23179dc3531a988e483b05d306539cca99f9a5c
RUN git clone $REPO --branch $VERSION --single-branch . && \
    git switch -c branch-$VERSION && \
    bash -c '[ "$(git rev-parse HEAD)" = "$COMMIT" ]' || (echo "Commit hash verification failed" && exit 1)

# Build the base-reth-node binary
RUN cargo build --bin base-reth-node --release

FROM ubuntu:22.04

# Install dependencies and cleanup
RUN apt-get update && \
    apt-get install -y jq curl supervisor && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /var/log/supervisor

WORKDIR /app

# Copy binaries and configuration
COPY --from=op /app/op-node/bin/op-node ./
COPY --from=reth /app/target/release/base-reth-node ./
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ./reth/reth-entrypoint ./execution-entrypoint
COPY op-node-entrypoint .

# Set execute permissions
RUN chmod +x ./execution-entrypoint ./op-node-entrypoint

# Add labels for better container identification
LABEL org.opencontainers.image.description="Base L2 node with reth execution client"
LABEL org.opencontainers.image.source="https://github.com/base/op-geth"

CMD ["/usr/bin/supervisord"]
