FROM golang:1.22 AS op

WORKDIR /app

ENV REPO=https://github.com/ethereum-optimism/optimism.git
ENV VERSION=v1.13.2
ENV COMMIT=c8b9f62736a7dad7e569719a84c406605f4472e6
RUN git clone $REPO --branch op-node/$VERSION --single-branch . && \
    git switch -c branch-$VERSION && \
    bash -c '[ "$(git rev-parse HEAD)" = "$COMMIT" ]'

# Install 'just'
RUN curl -sSfL 'https://just.systems/install.sh' | bash -s -- --to /usr/local/bin

RUN cd op-node && \
    make VERSION=$VERSION op-node

FROM rust:1.87 AS reth

WORKDIR /app

RUN apt-get update && apt-get -y upgrade && apt-get install -y git libclang-dev pkg-config curl build-essential

ENV REPO=https://github.com/paradigmxyz/reth.git
ENV VERSION=v1.4.3
ENV COMMIT=fe3653ffe602d4e85ad213e8bd9f06e7b710c0c5
RUN git clone $REPO --branch $VERSION --single-branch . && \
    git switch -c branch-$VERSION && \
    bash -c '[ "$(git rev-parse HEAD)" = "$COMMIT" ]'

RUN cargo build --bin op-reth --profile maxperf --manifest-path crates/optimism/bin/Cargo.toml

FROM rust:1.87 AS reth-base

WORKDIR /app

RUN apt-get update && apt-get -y upgrade && \
    apt-get install -y git libclang-dev pkg-config curl build-essential && \
    rm -rf /var/lib/apt/lists/*

ENV REPO=https://github.com/base/node-reth.git
ENV VERSION=v0.1.1
ENV COMMIT=cb55e69e3d88f8272ccc523c6c352d00b4ce2bf1
RUN git clone $REPO . && \
    git checkout tags/$VERSION && \
    bash -c '[ "$(git rev-parse HEAD)" = "$COMMIT" ]' || (echo "Commit hash verification failed" && exit 1)

RUN cargo build --bin base-reth-node --release

FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y jq curl supervisor && \
    rm -rf /var/lib/apt/lists
RUN mkdir -p /var/log/supervisor

WORKDIR /app

COPY --from=op /app/op-node/bin/op-node ./
COPY --from=reth /app/target/maxperf/op-reth ./
COPY --from=reth-base /app/target/release/base-reth-node ./
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ./reth/reth-entrypoint ./execution-entrypoint
COPY op-node-entrypoint .

CMD ["/usr/bin/supervisord"]
