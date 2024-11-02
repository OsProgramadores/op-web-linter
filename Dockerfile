# Defaults
ARG project_name="op-web-linter"
ARG project_title="OsProgramadores Web-based linter (op-web-linter)"
ARG project_author="paganini@paganini.net"
ARG project_source="https://github.com/osprogramadores/op-web-linter"
ARG project_user="op"
ARG project_uid=501

ARG home="/home/${project_user}"
ARG gopath="${home}/go"
ARG src_dir="${home}/${project_name}"

FROM alpine:3.17

LABEL org.opencontainers.image.title="${project_title}"
LABEL org.opencontainers.image.authors="${project_author}"
LABEL org.opencontainers.image.source="${project_source}"

# Pull from defaults.
ARG project_name
ARG project_user
ARG project_uid
ARG home
ARG gopath
ARG src_dir

# Fully static (as long we we don't need to link to C)
ENV CGO_ENABLED 0

# Directories & PATH
ENV HOME="${home}"
ENV GOPATH="${gopath}"
ENV SRC_DIR="${src_dir}"
ENV PATH="${PATH}:/usr/local/bin"

# Setup environment, install tools while under HOME since some
# tools (E.g. npm will use directories under the current location.)
WORKDIR ${home}

RUN apk add --no-cache ca-certificates clang15 clang15-extra-tools curl git git-crypt go indent make openjdk17 nodejs npm python3 py3-pylint && \
    adduser --uid ${project_uid} --home "${home}" --no-create-home --disabled-password ${project_user} && \
    mkdir -p "${gopath}" && \
    mkdir -p "${src_dir}" && \
    mkdir -p /usr/local/bin && \
    go install golang.org/x/lint/golint@latest && \
    cp "${gopath}/bin/golint" /usr/local/bin && \
    npm install --save-dev eslint-config-standard-with-typescript@23.0.0 eslint@8.24.0 && \
    curl -LJO "https://github.com/google/google-java-format/releases/download/v1.24.0/google-java-format-1.24.0-all-deps.jar"

# Copy repo contents, compile and install.
WORKDIR ${src_dir}
COPY . .
RUN go mod download && make install

# Default port.
EXPOSE 10000

USER ${project_uid}
ENTRYPOINT ["/usr/local/bin/op-web-linter"]
