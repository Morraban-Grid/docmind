#!/bin/bash
# Generate gRPC code from proto files
# This script generates Go and Python code from .proto definitions

set -e

echo "========================================="
echo "Generating gRPC Code"
echo "========================================="
echo ""

PROTO_DIR="proto/rag"
PROTO_FILE="$PROTO_DIR/rag.proto"

# Check if proto file exists
if [ ! -f "$PROTO_FILE" ]; then
    echo "❌ Proto file not found: $PROTO_FILE"
    exit 1
fi

# Generate Go code
echo "📦 Generating Go code..."
GO_OUT_DIR="$PROTO_DIR/gen/go"
mkdir -p "$GO_OUT_DIR"

protoc \
    --go_out="$GO_OUT_DIR" \
    --go_opt=paths=source_relative \
    --go-grpc_out="$GO_OUT_DIR" \
    --go-grpc_opt=paths=source_relative \
    "$PROTO_FILE"

echo "✅ Go code generated in $GO_OUT_DIR"

# Generate Python code
echo "📦 Generating Python code..."
PYTHON_OUT_DIR="$PROTO_DIR/gen/python"
mkdir -p "$PYTHON_OUT_DIR"

python -m grpc_tools.protoc \
    -I"$PROTO_DIR" \
    --python_out="$PYTHON_OUT_DIR" \
    --grpc_python_out="$PYTHON_OUT_DIR" \
    "$PROTO_FILE"

echo "✅ Python code generated in $PYTHON_OUT_DIR"

echo ""
echo "========================================="
echo "✅ Code generation complete!"
echo "========================================="
echo ""
