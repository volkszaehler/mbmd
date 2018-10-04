PWD := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BIN := $(PWD)/bin
BUILD := GOBIN=$(BIN) go install ./...
GOPATH := $(shell go env GOPATH)

all: build

build: assets binaries

binaries:
	@echo "Building for host platform"
	@$(BUILD)
	@echo "Created binaries:"
	@ls -1 bin

assets:
	@echo "Generating embedded assets"
	@$(GOPATH)/bin/embed http.go

release: test clean assets
	@./build.sh

test:
	@echo "Running testsuite"
	@go test

clean:
	@rm -rf bin/ pkg/ *.zip

dep:
	@echo "Installing embed tool"
	@go get -u github.com/aprice/embed/cmd/embed
	@echo "Installing vendor dependencies"
	@dep ensure

.PHONY: all build binaries assets release test clean dep
