package minio

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	// Client is the global MinIO client instance
	Client *minio.Client

	// Configuration
	MusicBucket string
	ImageBucket string
)

// Config holds MinIO configuration
type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	MusicBucket     string
	ImageBucket     string
}

// InitMinIO initializes the MinIO client with configuration from environment variables
func InitMinIO() error {
	config := Config{
		Endpoint:        getEnv("MINIO_ENDPOINT", "192.168.22.4:9100"),
		AccessKeyID:     getEnv("MINIO_ACCESS_KEY", "duylongadmin"),
		SecretAccessKey: getEnv("MINIO_SECRET_KEY", "duylongpass"),
		UseSSL:          getEnv("MINIO_USE_SSL", "false") == "true",
		MusicBucket:     getEnv("MINIO_MUSIC_BUCKET", "music"),
		ImageBucket:     getEnv("MINIO_IMAGE_BUCKET", "images"),
	}

	// Initialize MinIO client
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}

	Client = client
	MusicBucket = config.MusicBucket
	ImageBucket = config.ImageBucket

	// Test connection
	ctx := context.Background()
	_, err = Client.ListBuckets(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	log.Printf("✓ Connected to MinIO at %s", config.Endpoint)

	// Ensure buckets exist
	if err := ensureBucket(ctx, config.MusicBucket); err != nil {
		log.Printf("Warning: Music bucket '%s' check failed: %v", config.MusicBucket, err)
	}
	if err := ensureBucket(ctx, config.ImageBucket); err != nil {
		log.Printf("Warning: Image bucket '%s' check failed: %v", config.ImageBucket, err)
	}

	return nil
}

// ensureBucket checks if a bucket exists and creates it if it doesn't
func ensureBucket(ctx context.Context, bucketName string) error {
	exists, err := Client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf("Creating bucket: %s", bucketName)
		err = Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("✓ Created bucket: %s", bucketName)
	} else {
		log.Printf("✓ Bucket exists: %s", bucketName)
	}

	return nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetObject retrieves an object from MinIO
func GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

// StatObject gets object metadata
func StatObject(ctx context.Context, bucketName, objectName string) (minio.ObjectInfo, error) {
	return Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
}

// ListObjects lists all objects in a bucket
func ListObjects(ctx context.Context, bucketName string) <-chan minio.ObjectInfo {
	return Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})
}
