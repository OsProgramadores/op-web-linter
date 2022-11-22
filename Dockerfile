FROM golang:1-alpine AS builder

MAINTAINER Marco Paganini <paganini@paganini.net>

ARG PROJECT=op-web-linter

# Copy the repo contents into /tmp/build (inside the container)
WORKDIR /tmp/build/src/$PROJECT
COPY . .

# Fully static (as long we we don't need to link to C)
ENV CGO_ENABLED 0
ENV PATH="${PATH}:/usr/local/bin"

RUN apk add --no-cache ca-certificates git make nodejs npm && \
    export HOME="/build" && \
    export GOPATH="${HOME}" && \
    go mod download && \
    mkdir -p /usr/local/bin && \
    make install && \
    go get golang.org/x/lint/golint && \
    go install golang.org/x/lint/golint && \
    cp "${GOPATH}/bin/golint" /usr/local/bin && \
    mkdir -p /tmp/build/nodejs && \
    cd /tmp/build/nodejs && \
    npm install --save-dev eslint-config-standard-with-typescript@23.0.0 eslint@8.24.0

# Default port.
EXPOSE 10000

# ENTRYPOINT allows this image to be executed as a regular executable. Any
# arguments after the docker run are appended to the ENTRYPOINT command line.
USER ${UID}
ENTRYPOINT ["/usr/local/bin/op-web-linter"]

