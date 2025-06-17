FROM golang:1.23 AS op

RUN curl -sSfL 'https://just.systems/install.sh' | bash -s -- --to /usr/local/bin

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN . /tmp/versions.env && git clone $OP_NODE_REPO --branch $OP_NODE_TAG --single-branch . && \
    git switch -c branch-$OP_NODE_TAG && \
    bash -c '[ "$(git rev-parse HEAD)" = "$OP_NODE_COMMIT" ]'

RUN . /tmp/versions.env && cd op-node && \
    make VERSION=$OP_NODE_TAG op-node

FROM golang:1.23 AS geth

WORKDIR /app

COPY versions.env /tmp/versions.env

RUN . /tmp/versions.env && git clone $OP_GETH_REPO --branch $OP_GETH_TAG --single-branch . && git switch -c branch-$OP_GETH_TAG #&& bash -c '[ "$(git rev-parse HEAD)" = "$OP_GETH_COMMIT" ]'

#RUN . /tmp/versions.env && git clone $OP_GETH_REPO --branch $OP_GETH_TAG --single-branch . 
#RUN . /tmp/versions.env git switch -c branch-$OP_GETH_TAG 
#RUN . /tmp/versions.env bash -c '[ "$(git rev-parse HEAD)" = "$OP_GETH_COMMIT" ]'
