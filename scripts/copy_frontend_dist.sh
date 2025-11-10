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

# Optional: pre-compress common static assets to speed up production serving
# Create .gz files alongside originals (safe; will not remove originals)
if command -v gzip >/dev/null 2>&1; then
  echo "Pre-compressing assets (js/css/html/svg)..."
  find "$DEST_DIR" -type f \( -name "*.js" -o -name "*.css" -o -name "*.html" -o -name "*.svg" -o -name "*.txt" \) -print0 \
    | xargs -0 -r -n1 -P4 sh -c 'gzip -9 -c "$0" > "$0.gz"' 
  echo "Pre-compression done (created .gz alongside files)"
else
  echo "gzip not found; skipping pre-compression"
fi
