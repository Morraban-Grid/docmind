package infrastructure

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient wraps MinIO client operations
type MinIOClient struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		LogError("failed to create MinIO client", err, map[string]interface{}{
			"endpoint": endpoint,
		})
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Verify bucket exists, create if it doesn't
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		LogError("failed to check bucket existence", err, map[string]interface{}{
			"bucket": bucketName,
		})
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			LogError("failed to create bucket", err, map[string]interface{}{
				"bucket": bucketName,
			})
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		LogInfo("MinIO bucket created", map[string]interface{}{
			"bucket": bucketName,
		})
	}

	LogInfo("MinIO client initialized", map[string]interface{}{
		"endpoint": endpoint,
		"bucket":   bucketName,
	})

	return &MinIOClient{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// UploadFile uploads a file to MinIO and returns the storage path
func (m *MinIOClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, fileSize int64, contentType string) (string, error) {
	_, err := m.client.PutObject(
		ctx,
		m.bucketName,
		objectName,
		reader,
		fileSize,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		LogError("failed to upload file to MinIO", err, map[string]interface{}{
			"object_name": objectName,
			"file_size":   fileSize,
		})
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	storagePath := fmt.Sprintf("%s/%s", m.bucketName, objectName)

	LogInfo("file uploaded to MinIO", map[string]interface{}{
		"object_name":  objectName,
		"storage_path": storagePath,
		"file_size":    fileSize,
	})

	return storagePath, nil
}

// DownloadFile retrieves a file from MinIO
func (m *MinIOClient) DownloadFile(ctx context.Context, objectName string) (*minio.Object, error) {
	object, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		LogError("failed to download file from MinIO", err, map[string]interface{}{
			"object_name": objectName,
		})
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	// Verify object exists by checking stat
	_, err = object.Stat()
	if err != nil {
		LogError("file not found in MinIO", err, map[string]interface{}{
			"object_name": objectName,
		})
		return nil, fmt.Errorf("file not found: %w", err)
	}

	LogInfo("file downloaded from MinIO", map[string]interface{}{
		"object_name": objectName,
	})

	return object, nil
}

// DeleteFile removes a file from MinIO
func (m *MinIOClient) DeleteFile(ctx context.Context, objectName string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		LogError("failed to delete file from MinIO", err, map[string]interface{}{
			"object_name": objectName,
		})
		return fmt.Errorf("failed to delete file: %w", err)
	}

	LogInfo("file deleted from MinIO", map[string]interface{}{
		"object_name": objectName,
	})

	return nil
}

// GetPresignedURL generates a presigned URL for temporary file access
func (m *MinIOClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expiry, nil)
	if err != nil {
		LogError("failed to generate presigned URL", err, map[string]interface{}{
			"object_name": objectName,
		})
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}
