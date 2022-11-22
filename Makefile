## This file is part of op-web-linter.
## See github.com/osprogramadores/op-web-linter for licensing and details.

.PHONY: arch clean docker install

BIN := op-web-linter
BINDIR := /usr/local/bin
ARCHDIR := arch
SRC := $(wildcard *.go) $(wildcard common/*.go) $(wildcard handlers/*.go) $(wildcard lang/*.go) templates/form.tmpl
GIT_TAG := $(shell git describe --always --tags)

# Default target
${BIN}: Makefile ${SRC}
	CGO_ENABLED=0 go build -v -ldflags "-X main.BuildVersion=${GIT_TAG}" -o "${BIN}"

clean:
	rm -f "${BIN}"
	rm -f "docs/${BIN}.1"
	rm -rf "${ARCHDIR}"

install: ${BIN}
	install -m 755 "${BIN}" "${BINDIR}"

docker: ${BIN}
	docker build -t "${BIN}:latest" .

# Creates cross-compiled tarred versions (for releases).
arch: Makefile ${SRC}
	for ga in "linux/amd64" "linux/386" "linux/arm" "linux/arm64" "linux/mips" "linux/mipsle"; do \
	  export goos="$${ga%/*}"; \
	  export goarch="$${ga#*/}"; \
	  dst="./${ARCHDIR}/$${goos}-$${goarch}"; \
	  mkdir -p "$${dst}"; \
	  echo "=== Building $${goos}/$${goarch} ==="; \
	  go build -v -ldflags "-X main.Build=${GIT_TAG}" -o "$${dst}/${BIN}"; \
	  [ -s LICENSE ] && install -m 644 LICENSE "$${dst}"; \
	  [ -s README.md ] && install -m 644 README.md "$${dst}"; \
	  [ -s docs/${BIN}.1 ] && install -m 644 docs/${BIN}.1 "$${dst}"; \
	  tar -C "${ARCHDIR}" -zcvf "${ARCHDIR}/${BIN}-$${goos}-$${goarch}.tar.gz" "$${dst##*/}"; \
	done
