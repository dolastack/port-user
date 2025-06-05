#!/bin/bash

set -e

BINARY_NAME="port-user"
OUTPUT_DIR="dist"

rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

build() {
    local os="$1"
    local arch="$2"
    local ext="${3:-}"
    local target="$OUTPUT_DIR/$BINARY_NAME-$os-$arch$ext"

    echo "Building for $os/$arch..."
    GOOS="$os" GOARCH="$arch" go build -o "$target" main.go
    chmod +x "$target"
    echo "Saved to: $target"
}

# Linux
build linux amd64
build linux arm64

# macOS
build darwin amd64
build darwin arm64

# Windows
build windows amd64 .exe

echo ""
echo "✅ All builds completed successfully!"