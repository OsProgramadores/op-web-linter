FROM alpine:3.17 as builder

MAINTAINER Marco Paganini <paganini@paganini.net>

ARG PROJECT=op-web-linter
ARG PROJECT_UID=501
ARG PROJECT_USER=op

# Fully static (as long we we don't need to link to C)
ENV CGO_ENABLED 0

# Directories & PATH
ENV HOME="/home/$PROJECT_USER"
ENV GOPATH="${HOME}/go"
ENV SRC_DIR="${HOME}/${PROJECT}"
ENV PATH="${PATH}:/usr/local/bin"

WORKDIR $HOME

RUN apk add --no-cache ca-certificates clang15 clang15-extra-tools curl git git-crypt go indent make openjdk17 nodejs npm python3 py3-pylint && \
    adduser --uid $PROJECT_UID --home /tmp --no-create-home --disabled-password $PROJECT_USER && \
    mkdir -p "${HOME}" && \
    mkdir -p "${GOPATH}" && \
    mkdir -p "${SRC_DIR}" && \
    mkdir -p /usr/local/bin && \
    go install golang.org/x/lint/golint@latest && \
    cp "${GOPATH}/bin/golint" /usr/local/bin && \
    npm install --save-dev eslint-config-standard-with-typescript@23.0.0 eslint@8.24.0 && \
    curl -LJO "https://github.com/google/google-java-format/releases/download/v1.15.0/google-java-format-1.15.0-all-deps.jar"

WORKDIR $SRC_DIR
COPY . .

RUN go mod download && \
    make install

# Default port.
EXPOSE 10000

# ENTRYPOINT allows this image to be executed as a regular executable. Any
# arguments after the docker run are appended to the ENTRYPOINT command line.
USER ${PROJECT_UID}
ENTRYPOINT ["/usr/local/bin/op-web-linter"]
