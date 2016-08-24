SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=otp
BINARY_OUT=bin

# These are the values we want to pass for Version and BuildTime
VERSION=0.0.2
BUILD_TIME=$(shell date +%FT%T%z)
BUILD_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_ARCH=$(shell uname -m)
BUILD_OS=$(shell uname | tr '[:upper:]' '[:lower:]')

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
		mkdir -p ${BINARY_OUT} && \
    go build -o ${BINARY_OUT}/${BINARY} ./cmd/otp/main.go

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	rm -rf ${BINARY_OUT} _release

.PHONY: release
release:
	[ ! -d bin ] && echo "Run make first." && exit 1;\
	mkdir -p _release && \
	tar czvf _release/dev-tools-${VERSION}-${BUILD_COMMIT}-${BUILD_ARCH}-${BUILD_OS}.tar.gz -C bin .
