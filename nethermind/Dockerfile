FROM golang:1.23 AS op

RUN curl -sSfL 'https://just.systems/install.sh' | bash -s -- --to /usr/local/bin

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN . /tmp/versions.env && git clone $OP_NODE_REPO --branch $OP_NODE_TAG --single-branch . && \
    git switch -c branch-$OP_NODE_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$OP_NODE_COMMIT" ]'

RUN . /tmp/versions.env && cd op-node && \
    just VERSION=$OP_NODE_TAG op-node

FROM mcr.microsoft.com/dotnet/sdk:9.0-noble AS build

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

FROM mcr.microsoft.com/dotnet/aspnet:9.0-noble

RUN apt-get update && \
    apt-get install -y jq curl supervisor && \
    rm -rf /var/lib/apt/lists

RUN mkdir -p /var/log/supervisor

WORKDIR /app

COPY --from=build /publish ./
COPY --from=op /app/op-node/bin/op-node ./
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY nethermind/nethermind-entrypoint ./execution-entrypoint
COPY op-node-entrypoint .

CMD ["/usr/bin/supervisord"]
