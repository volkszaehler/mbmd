.PHONY: default clean install lint test build binaries publish-images test-release release

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_TAGS := -tags=release
LD_FLAGS := -X "github.com/volkszaehler/mbmd/server.Version=${VERSION}" -X "github.com/volkszaehler/mbmd/server.Commit=${SHA}"
BUILD_ARGS := -ldflags='$(LD_FLAGS)'

default: clean install lint test build

clean:
	rm -rf dist/

install:
	go install github.com/alvaroloes/enumer

lint:
	golangci-lint run --timeout 5m

test:
	@echo "Running testsuite"
	go test ./...

build:
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v $(BUILD_TAGS) $(BUILD_ARGS)

publish-images:
	@echo Version: $(VERSION) $(BUILD_DATE)
	seihon publish -v "$(TAG_NAME)" -v "latest" --image-name volkszaehler/mbmd --base-runtime-image alpine --dry-run=false

test-release:
	goreleaser --snapshot --skip-publish --rm-dist

release:
	goreleaser --rm-dist
