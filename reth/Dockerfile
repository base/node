FROM golang:1.24 AS op

RUN curl -sSfL 'https://just.systems/install.sh' | bash -s -- --to /usr/local/bin

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN . /tmp/versions.env && git clone $OP_NODE_REPO --branch $OP_NODE_TAG --single-branch . && \
    git switch -c branch-$OP_NODE_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$OP_NODE_COMMIT" ]'

RUN . /tmp/versions.env && cd op-node && \
    make VERSION=$OP_NODE_TAG op-node

FROM rust:1.88 AS reth-base

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN apt-get update && apt-get -y upgrade && \
    apt-get install -y git libclang-dev pkg-config curl build-essential && \
    rm -rf /var/lib/apt/lists/*

RUN . /tmp/versions.env && git clone $BASE_RETH_NODE_REPO . && \
    git checkout tags/$BASE_RETH_NODE_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$BASE_RETH_NODE_COMMIT" ]' || (echo "Commit hash verification failed" && exit 1)

RUN cargo build --bin base-reth-node --profile maxperf

FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y jq curl supervisor && \
    rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/log/supervisor

WORKDIR /app

COPY --from=op /app/op-node/bin/op-node ./
COPY --from=reth-base /app/target/maxperf/base-reth-node ./
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ./reth/reth-entrypoint ./execution-entrypoint
COPY op-node-entrypoint .

CMD ["/usr/bin/supervisord"]
