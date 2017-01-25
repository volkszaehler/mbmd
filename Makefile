all: build

build:
	@echo "Building for host platform"
	@gb build all
	@echo "Building binary for Raspberry Pi"
	@GOOS=linux GOARCH=arm GOARM=5 gb build all
	@echo "Created binaries:"
	@ls -1 bin

clean:
	@rm -rf bin/ pkg/

dep:
	@echo "Installing GB build tool"
	@go get github.com/constabulary/gb/...

.PHONY: all build clean
