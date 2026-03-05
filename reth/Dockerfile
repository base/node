FROM rust:1.88 AS base

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN apt-get update && apt-get -y upgrade && \
    apt-get install -y git libclang-dev pkg-config curl build-essential && \
    rm -rf /var/lib/apt/lists/*

RUN . /tmp/versions.env && git clone $BASE_NODE_REPO . && \
    git checkout tags/$BASE_NODE_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$BASE_NODE_COMMIT" ]' || (echo "Commit hash verification failed" && exit 1)

RUN cargo build --bin base-reth-node --bin base-consensus --profile maxperf

FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y jq curl supervisor && \
    rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/log/supervisor

WORKDIR /app

COPY --from=base /app/target/maxperf/base-consensus ./
COPY --from=base /app/target/maxperf/base-reth-node ./
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ./reth/reth-entrypoint ./execution-entrypoint
COPY base-consensus-entrypoint .

CMD ["/usr/bin/supervisord"]
