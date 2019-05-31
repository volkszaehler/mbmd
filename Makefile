PWD := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BIN := $(PWD)/bin
BUILD := GO111MODULE=on GOBIN=$(BIN) go install ./...

all: build




release: test clean assets

test:
	@echo "Running testsuite"
	GO111MODULE=on go test ./...

clean:
	rm -rf bin/ pkg/ *.zip
build: assets binaries

assets:
	@echo "Generating embedded assets"
	GO111MODULE=on go generate ./...

binaries:
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "github.com/gonium/gosdm630.Version=${VERSION}" -X "github.com/gonium/gosdm630.Commit=${SHA}"' ./cmd/sdm

.PHONY: all build binaries assets release test clean
