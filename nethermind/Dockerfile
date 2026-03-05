FROM rust:1.88 AS base-consensus

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN apt-get update && apt-get -y upgrade && \
    apt-get install -y git libclang-dev pkg-config curl build-essential && \
    rm -rf /var/lib/apt/lists/*

RUN . /tmp/versions.env && git clone $BASE_NODE_REPO . && \
    git checkout $BASE_NODE_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$BASE_NODE_COMMIT" ]' || (echo "Commit hash verification failed" && exit 1)

RUN cargo build --bin base-consensus --profile maxperf

FROM mcr.microsoft.com/dotnet/sdk:10.0-noble AS build

ARG BUILD_CONFIG=release
ARG TARGETARCH

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN . /tmp/versions.env && git clone $NETHERMIND_REPO --branch $NETHERMIND_TAG --single-branch . && \
    git switch -c $NETHERMIND_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$NETHERMIND_COMMIT" ]'
    
RUN TARGETARCH=${TARGETARCH#linux/} && \
    arch=$([ "$TARGETARCH" = "amd64" ] && echo "x64" || echo "$TARGETARCH") && \
    echo "Using architecture: $arch" && \
    dotnet publish src/Nethermind/Nethermind.Runner -c $BUILD_CONFIG -a $arch -o /publish --sc false

FROM mcr.microsoft.com/dotnet/aspnet:10.0-noble

RUN apt-get update && \
    apt-get install -y jq curl supervisor && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/log/supervisor

WORKDIR /app

COPY --from=build /publish ./
COPY --from=base-consensus /app/target/maxperf/base-consensus ./
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY nethermind/nethermind-entrypoint ./execution-entrypoint
COPY base-consensus-entrypoint .

CMD ["/usr/bin/supervisord"]
