#!/bin/bash


# Parameters
WORK_DIR=$(pwd)
SOURCE_DIR="$WORK_DIR/wishlist"
OUTPUT_BINARY="/app/wishlist"
BIN_DIR="./bin"


# Run the container for building
docker run --rm -it \
  -v "$SOURCE_DIR:/app" \
  -w /app \
  golang:1.23 \
  /bin/bash -c "set -x; go mod download; go build -o $OUTPUT_BINARY ."

# Check if the directory exists
if [ ! -d "$DIR_BIN_DIRPATH" ]; then
  mkdir -p "$BIN_DIR"
  echo "Directory '$BIN_DIR' created."
else
  echo "Directory '$BIN_DIR' already exists."
fi

# Check build success
if [ $? -eq 0 ]; then
  echo "Build completed successfully."
  mv "$SOURCE_DIR/wishlist" "$BIN_DIR/"
  echo "The binary file has been moved to: $BIN_DIR/wishlist"
else
  echo "Build failed with an error."
  exit 1
fi

