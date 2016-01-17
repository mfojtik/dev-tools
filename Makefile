SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=otp
BINARY_OUT=bin

# These are the values we want to pass for Version and BuildTime
VERSION=0.0.1
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
	if [ -d ${BINARY_OUT} ] ; then rm -rf ${BINARY_OUT} ; fi

.PHONY: release
release:
	tar czvf dev-tools-${VERSION}-${BUILD_COMMIT}-${BUILD_ARCH}-${BUILD_OS}.tar.gz -C bin .
