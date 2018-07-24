#!/bin/bash

RELEASE=$(pwd)/release
BIN=$(pwd)/bin
CMDS=$(go list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...)

mkdir -p $RELEASE
mkdir -p $BIN
if ls *.zip 1>/dev/null 2>&1; then rm *.zip; fi

function build {
	GOOS=$1
	GOARCH=$2
	GOARM=$3

	cd $BIN
	if [ ! -z "$(ls -A .)" ]; then rm *; fi
	for i in $CMDS; do
		# echo GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM go build $i
		GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM go build $i
	done
	cd ..

	zip $RELEASE/sdm630-$GOOS-$GOARCH $BIN/*
}

echo "Building for ..."
echo "... Linux/32bit"
build linux 386
echo "... Linux/64bit"
build linux amd64
echo "... Raspberry Pi/Linux"
build linux arm 5
echo "... Mac OS/64bit"
build darwin amd64
echo "... Windows/32bit"
build windows 386
echo "... Windows/64bit"
build windows amd64
