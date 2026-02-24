#!/usr/bin/env bash
set -e

# ---------------------------------------------
# Script to generate protobuf code for docmind
# ---------------------------------------------

# Base directory of proto files (absolute path)
BASE_DIR="$(cd "$(dirname "$0")/../proto/rag" && pwd)"
GEN_GO_DIR="$BASE_DIR/gen/go"
GEN_PY_DIR="$BASE_DIR/gen/python"

# Create directories if they do not exist
mkdir -p "$GEN_GO_DIR"
mkdir -p "$GEN_PY_DIR"

echo "Generating Go protobuf files..."
protoc \
  --proto_path="$BASE_DIR" \
  --go_out="$GEN_GO_DIR" \
  --go-grpc_out="$GEN_GO_DIR" \
  "$BASE_DIR/rag.proto"

echo "Generating Python protobuf files..."
python -m grpc_tools.protoc \
  --proto_path="$BASE_DIR" \
  --python_out="$GEN_PY_DIR" \
  --grpc_python_out="$GEN_PY_DIR" \
  "$BASE_DIR/rag.proto"

echo "Protobuf generation completed successfully."