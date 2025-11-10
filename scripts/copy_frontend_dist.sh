#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/copy_frontend_dist.sh /path/to/frontend/project
# If no arg provided, reads FRONT_SRC env or defaults to ../Downloads/carbon-scan-now-main

SRC=${1:-${FRONT_SRC:-"../Downloads/carbon-scan-now-main"}}
DEST_DIR="$(cd "$(dirname "$0")/.." && pwd)/frontend/dist"

echo "Frontend source: $SRC"
echo "Destination: $DEST_DIR"

if [ ! -d "$SRC" ]; then
  echo "Source directory not found: $SRC"
  exit 2
fi

# Install and build
cd "$SRC"
if [ ! -d node_modules ]; then
  echo "Installing frontend dependencies..."
  npm install
fi

echo "Building frontend..."
npm run build

# Find build dir (Vite default is dist)
if [ ! -d dist ]; then
  echo "Build output 'dist' not found in $SRC"
  exit 3
fi

# Prepare destination
mkdir -p "$DEST_DIR"
rm -rf "$DEST_DIR"/*

echo "Copying dist to $DEST_DIR"
cp -r dist/* "$DEST_DIR/"

echo "Done. Frontend copied to $DEST_DIR"
