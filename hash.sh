#!/bin/bash

HASH=$(git rev-parse HEAD | head -c8)
TAG=$(git describe --abbrev=0 --tags)

cat << EOF > version.go
package sdm630

const (
	TAG = "$TAG"
	HASH = "$HASH"
)
EOF
