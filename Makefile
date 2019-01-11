PWD := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BIN := $(PWD)/bin
BUILD := GO111MODULE=on GOBIN=$(BIN) go install ./...

all: build

build: assets binaries

binaries:
	@echo "Building for host platform"
	$(BUILD)
	@echo "Created binaries:"
	@ls -1 bin

assets:
	./hash.sh
	@echo "Generating embedded assets"
	GO111MODULE=on go generate ./...

release: test clean assets
	./build.sh

test:
	@echo "Running testsuite"
	GO111MODULE=on go test ./...

clean:
	rm -rf bin/ pkg/ *.zip

.PHONY: all build binaries assets release test clean
