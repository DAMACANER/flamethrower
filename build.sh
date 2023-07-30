#!/bin/bash

# Define the name of your application
APP_NAME="flamethrower-tui"

# Define the main .go file
MAIN_FILE="main.go"

# Get the current version from git
VERSION=$(git describe --tags)

# Function to build for a specific OS/ARCH
build() {
  echo "Building for $1 $2..."
  env GOOS=$1 GOARCH=$2 go build -ldflags "-X main.version=$VERSION" -o build/$APP_NAME-$1-$2 $MAIN_FILE
}

mkdir -p build

# Build for each OS/ARCH
build linux amd64
build windows amd64
build darwin amd64

echo "Build complete"
