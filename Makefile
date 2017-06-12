all: build

build:
	@echo "Building for host platform"
	@gb build all
	@echo "Building binary for Raspberry Pi"
	@GOOS=linux GOARCH=arm GOARM=5 gb build all
	@echo "Created binaries:"
	@ls -1 bin

release-build: test clean
	@echo "Building binaries..."
	@echo "... for Linux/32bit"
	@GOOS=linux GOARCH=386 gb build all
	@echo "... for Linux/64bit"
	@GOOS=linux GOARCH=amd64 gb build all
	@echo "... for Raspberry Pi/Linux"
	@GOOS=linux GOARCH=arm GOARM=5 gb build all
	@echo "... for Mac OS/64bit"
	@GOOS=darwin GOARCH=amd64 gb build all
	@echo "... for Windows/32bit"
	@GOOS=windows GOARCH=386 gb build all
	@echo "... for Linux/64bit"
	@GOOS=windows GOARCH=amd64 gb build all
	@echo
	@echo "Created binaries:"
	@ls -1 bin

release: release-build
	@echo "Starting packaging"
	@echo "... for Linux"
	@zip sdm630-linux-386 bin/*-linux-386*
	@zip sdm630-linux-amd64 bin/*-linux-amd64
	@zip sdm630-linux-arm bin/*-linux-arm*
	@echo "... for Mac OS"
	@zip sdm630-darwin-amd64 bin/*-darwin-amd64
	@echo "... for Windows"
	@zip sdm630-windows-386 bin/*-windows-386*

test:
	@echo "Running testsuite"
	@gb test

clean:
	@rm -rf bin/ pkg/ *.zip

dep:
	@echo "Installing GB build tool"
	@go get github.com/constabulary/gb/...

.PHONY: all build clean
