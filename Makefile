.PHONY: default clean checks lint test build assets binaries publish-images test-release

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

default: clean checks test build

clean:
	rm -rf dist/ *.zip

checks: assets lint

lint:
	golangci-lint run

test:
	@echo "Running testsuite"
	go test ./...

build: assets binaries

assets:
	@echo "Generating embedded assets"
	go generate ./...

binaries:
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "github.com/volkszaehler/mbmd/server.Version=${VERSION}" -X "github.com/volkszaehler/mbmd/server.Commit=${SHA}"' ./cmd/mbmd

publish-images:
	@echo Version: $(VERSION) $(BUILD_DATE)
	seihon publish -v "$(TAG_NAME)" -v "latest" --image-name volkszaehler/mbmd --base-runtime-image alpine --dry-run=false

test-release:
	goreleaser --snapshot --skip-publish --rm-dist
