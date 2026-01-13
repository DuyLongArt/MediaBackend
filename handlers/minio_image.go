package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	minioClient "MediaBackend/minio"
)

// StreamMinIOImage handles streaming image files from MinIO with caching support
func StreamMinIOImage(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	filename := strings.TrimPrefix(r.URL.Path, "/api/images/")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Get object info for metadata
	objectInfo, err := minioClient.StatObject(ctx, minioClient.ImageBucket, filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("Error getting object info for %s: %v", filename, err)
		return
	}

	// Set content type based on file extension
	contentType := getImageContentType(filename)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(objectInfo.Size, 10))

	// Add caching headers
	w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	w.Header().Set("Last-Modified", objectInfo.LastModified.UTC().Format(http.TimeFormat))

	// Use MinIO ETag for cache validation
	etag := objectInfo.ETag
	w.Header().Set("ETag", etag)

	// Check If-None-Match header for ETag validation
	if match := r.Header.Get("If-None-Match"); match != "" {
		if match == etag {
			w.WriteHeader(http.StatusNotModified)
			log.Printf("MinIO image not modified (ETag match): %s", filename)
			return
		}
	}

	// Check If-Modified-Since header
	if modifiedSince := r.Header.Get("If-Modified-Since"); modifiedSince != "" {
		if t, err := time.Parse(http.TimeFormat, modifiedSince); err == nil {
			if objectInfo.LastModified.Before(t.Add(1 * time.Second)) {
				w.WriteHeader(http.StatusNotModified)
				log.Printf("MinIO image not modified (time-based): %s", filename)
				return
			}
		}
	}

	// Get the object
	object, err := minioClient.GetObject(ctx, minioClient.ImageBucket, filename)
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		log.Printf("Error getting object: %v", err)
		return
	}
	defer object.Close()

	// Serve the file
	w.WriteHeader(http.StatusOK)
	io.Copy(w, object)
	log.Printf("Served MinIO image: %s (%d bytes)", filename, objectInfo.Size)
}

// ListMinIOImages returns a list of available image files in MinIO
func ListMinIOImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var files []MediaFile
	objectCh := minioClient.ListObjects(ctx, minioClient.ImageBucket)

	for object := range objectCh {
		if object.Err != nil {
			log.Printf("Error listing objects: %v", object.Err)
			continue
		}

		name := filepath.Base(object.Key)
		files = append(files, MediaFile{
			Name:        name,
			Size:        object.Size,
			Path:        object.Key,
			Url:         "/gomedia/api/images/" + name,
			ContentType: getImageContentType(name),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"files": %s}`, formatFileList(files))))
}
