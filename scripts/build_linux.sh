#!/usr/bin/env bash
set -euo pipefail

# Build Linux binary
mkdir -p dist
GOOS=linux GOARCH=amd64 go build -o dist/server-linux-amd64 ./

echo "Built dist/server-linux-amd64"
