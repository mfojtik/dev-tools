SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=otp
BINARY_OUT="bin"

# These are the values we want to pass for Version and BuildTime
VERSION=0.0.1
BUILD_TIME=`date +%FT%T%z`
BUILD_COMMIT=$(shell git rev-parse --short HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X github.com/mfojtik/dev-tools/core.Version=${VERSION}-${BUILD_COMMIT} -X github.com/mfojtik/dev-tools/core.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
		mkdir -p ${BINARY_OUT} && \
    go build ${LDFLAGS} \
		-o ${BINARY_OUT}/${BINARY} ./cmd/otp/main.go

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -d ${BINARY_OUT} ] ; then rm -rf ${BINARY_OUT} ; fi
