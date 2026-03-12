#!/bin/bash
# MinIO bucket initialization script
# This script creates the required bucket for document storage

set -e

echo "Waiting for MinIO to be ready..."
sleep 5

# Install MinIO client if not present
if ! command -v mc &> /dev/null; then
    echo "Installing MinIO client..."
    wget https://dl.min.io/client/mc/release/linux-amd64/mc -O /usr/local/bin/mc
    chmod +x /usr/local/bin/mc
fi

# Configure MinIO client
echo "Configuring MinIO client..."
mc alias set docmind http://minio:9000 ${MINIO_ROOT_USER:-minioadmin} ${MINIO_ROOT_PASSWORD:-minioadmin}

# Create bucket if it doesn't exist
BUCKET_NAME=${MINIO_BUCKET:-docmind-documents}
echo "Creating bucket: $BUCKET_NAME"

if mc ls docmind/$BUCKET_NAME > /dev/null 2>&1; then
    echo "Bucket $BUCKET_NAME already exists"
else
    mc mb docmind/$BUCKET_NAME
    echo "Bucket $BUCKET_NAME created successfully"
fi

# Set bucket policy to private
echo "Setting bucket policy to private..."
mc anonymous set none docmind/$BUCKET_NAME

echo "MinIO initialization complete!"
