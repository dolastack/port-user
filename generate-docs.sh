#!/bin/bash

set -e

BINARY="port-user"

# Build binary
echo "Building $BINARY..."
go build -o "$BINARY" main.go

# Generate docs
echo "Generating man pages..."
./$BINARY gen-docs

# Clean up
rm -f "$BINARY"

echo "Done! Man pages are in ./docs/"