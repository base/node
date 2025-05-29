FROM golang:1.23 AS op

WORKDIR /app

ENV REPO=https://github.com/ethereum-optimism/optimism.git
ENV VERSION=v1.13.3
ENV COMMIT=b1e7c63bb2ffea46771c302bcb05f72ba1a7bf61
RUN git clone $REPO --branch op-node/$VERSION --single-branch . && \
    git switch -c branch-$VERSION && \
    bash -c '[ "$(git rev-parse HEAD)" = "$COMMIT" ]'

# Install 'just'
RUN curl -sSfL 'https://just.systems/install.sh' | bash -s -- --to /usr/local/bin

RUN cd op-node && \
    make VERSION=$VERSION op-node

FROM golang:1.23 AS geth

WORKDIR /app

RUN apt-get update && \
    apt-get install -y \
    build-essential

ENV REPO=https://github.com/ethereum-optimism/op-geth.git
ENV VERSION=v1.101511.0
ENV COMMIT=68075997f33907401a93216aa426514c5ddc8870
RUN git clone $REPO --branch $VERSION --single-branch . && \
    git switch -c branch-$VERSION && \
    bash -c '[ "$(git rev-parse HEAD)" = "$COMMIT" ]'

RUN go run build/ci.go install -static ./cmd/geth

FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install -y jq curl supervisor && \
    rm -rf /var/lib/apt/lists
RUN mkdir -p /var/log/supervisor

WORKDIR /app

COPY --from=op /app/op-node/bin/op-node ./
COPY --from=geth /app/build/bin/geth ./
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY geth/geth-entrypoint ./execution-entrypoint
COPY op-node-entrypoint .

CMD ["/usr/bin/supervisord"]
