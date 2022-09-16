.PHONY: default clean docs porcelain install assets lint test build publish-images test-release release

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_TAGS := -tags=release
MODULE := github.com/volkszaehler/mbmd
LD_FLAGS := -X "${MODULE}/server.Version=${VERSION}" -X "${MODULE}/server.Commit=${SHA}"
BUILD_ARGS := -ldflags='$(LD_FLAGS)'

default: clean install lint test build

clean:
	rm -rf dist/

docs:
	go run $(MODULE) doc

porcelain:
	gofmt -w -l $$(find . -name '*.go')
	go mod tidy
	test -z "$$(git status --porcelain)" || (git status; git diff; false)

install:
	go install $$(go list -f '{{join .Imports " "}}' tools.go)

assets:
	go generate ./...

lint:
	golangci-lint run

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
