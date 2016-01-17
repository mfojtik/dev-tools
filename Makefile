SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=otp
BINARY_OUT="bin"

# These are the values we want to pass for Version and BuildTime
VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X github.com/mfojtik/openshift-dev-tools/core.Version=${VERSION} -X github.com/mfojtik/openshift-dev-tools/core.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
		mkdir -p ${BINARY_OUT} && \
    go build ${LDFLAGS} \
		-o ${BINARY_OUT}/${BINARY} ./cmd/openshift-tag-pr/main.go

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -d ${BINARY_OUT} ] ; then rm -rf ${BINARY_OUT} ; fi
