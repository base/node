ARG RUST_VERSION=1.93

FROM golang:1.24 AS op

RUN curl -sSfL 'https://just.systems/install.sh' | bash -s -- --to /usr/local/bin

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN . /tmp/versions.env && git clone $OP_NODE_REPO --branch $OP_NODE_TAG --single-branch . && \
    git switch -c branch-$OP_NODE_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$OP_NODE_COMMIT" ]'

RUN . /tmp/versions.env && cd op-node && \
    make VERSION=$OP_NODE_TAG op-node

FROM public.ecr.aws/docker/library/rust:${RUST_VERSION}-trixie AS rust-builder-base

WORKDIR /app

ARG MOLD_VERSION=2.40.4
ARG MOLD_SHA256_AARCH64=c799b9ccae8728793da2186718fbe53b76400a9da396184fac0c64aa3298ec37
ARG MOLD_SHA256_ARM=d82792748a81202423ddd2496fc8719404fe694493abdef691cc080392ee44bf
ARG MOLD_SHA256_X86_64=4c999e19ffa31afa5aa429c679b665d5e2ca5a6b6832ad4b79668e8dcf3d8ec1

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      git libclang-dev pkg-config curl build-essential cmake && \
    rm -rf /var/lib/apt/lists/*

RUN set -eux; \
    case "$(uname -m)" in \
      x86_64) MOLD_ARCH=x86_64; MOLD_SHA256="${MOLD_SHA256_X86_64}" ;; \
      aarch64|arm64) MOLD_ARCH=aarch64; MOLD_SHA256="${MOLD_SHA256_AARCH64}" ;; \
      armv7l|armv6l) MOLD_ARCH=arm; MOLD_SHA256="${MOLD_SHA256_ARM}" ;; \
      *) echo "unsupported architecture: $(uname -m)" >&2; exit 1 ;; \
    esac; \
    curl -fsSL "https://github.com/rui314/mold/releases/download/v${MOLD_VERSION}/mold-${MOLD_VERSION}-${MOLD_ARCH}-linux.tar.gz" -o /tmp/mold.tar.gz; \
    echo "${MOLD_SHA256}  /tmp/mold.tar.gz" | sha256sum -c -; \
    tar -xzf /tmp/mold.tar.gz -C /tmp; \
    cp /tmp/mold-${MOLD_VERSION}-${MOLD_ARCH}-linux/bin/* /usr/local/bin/; \
    rm -rf /tmp/mold*

FROM rust-builder-base AS reth-base

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN . /tmp/versions.env && git clone $BASE_RETH_NODE_REPO . && \
    git checkout tags/$BASE_RETH_NODE_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$BASE_RETH_NODE_COMMIT" ]' || (echo "Commit hash verification failed" && exit 1)

RUN cargo build --bin base-reth-node --bin base-consensus --profile maxperf

FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y jq curl supervisor && \
    rm -rf /var/lib/apt/lists/*
RUN mkdir -p /var/log/supervisor

WORKDIR /app

COPY --from=op /app/op-node/bin/op-node ./
COPY --from=reth-base /app/target/maxperf/base-consensus ./
COPY --from=reth-base /app/target/maxperf/base-reth-node ./
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ./reth/reth-entrypoint ./execution-entrypoint
COPY helpers.sh .
COPY op-node-entrypoint .
COPY base-consensus-entrypoint .
COPY consensus-entrypoint .

CMD ["/usr/bin/supervisord"]
