FROM golang:1-alpine AS builder

MAINTAINER Marco Paganini <paganini@paganini.net>

ARG PROJECT=op-web-linter
ARG BUILD_DIR=/tmp/build/src/$PROJECT

# Fully static (as long we we don't need to link to C)
ENV CGO_ENABLED 0
ENV PATH="${PATH}:/usr/local/bin"

RUN apk add --no-cache ca-certificates git make nodejs npm && \
    export HOME="/tmp/build" && \
    export GOPATH="/tmp/build" && \
    mkdir -p /usr/local/bin && \
    go install golang.org/x/lint/golint@latest && \
    cp "${GOPATH}/bin/golint" /usr/local/bin && \
    mkdir -p /tmp/build/nodejs && \
    cd /tmp/build/nodejs && \
    npm install --save-dev eslint-config-standard-with-typescript@23.0.0 eslint@8.24.0

WORKDIR $BUILD_DIR
COPY . .

RUN cd "${BUILD_DIR}" && \
    go mod download && \
    make install

# Default port.
EXPOSE 10000

# ENTRYPOINT allows this image to be executed as a regular executable. Any
# arguments after the docker run are appended to the ENTRYPOINT command line.
USER ${UID}
ENTRYPOINT ["/usr/local/bin/op-web-linter"]

