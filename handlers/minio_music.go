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

	minioClient "MediaBackend/minio"

	"github.com/minio/minio-go/v7"
)

// StreamMinIOMusic handles streaming music files from MinIO with HTTP range request support
func StreamMinIOMusic(w http.ResponseWriter, r *http.Request) {
	// Extract filename from URL path
	filename := strings.TrimPrefix(r.URL.Path, "/api/music/")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Get object info for metadata
	objectInfo, err := minioClient.StatObject(ctx, minioClient.MusicBucket, filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("Error getting object info for %s: %v", filename, err)
		return
	}

	fileSize := objectInfo.Size

	// Set content type based on file extension
	contentType := getContentType(filename)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("ETag", objectInfo.ETag)
	w.Header().Set("Last-Modified", objectInfo.LastModified.UTC().Format(http.TimeFormat))

	// Handle range requests for streaming
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		// No range request, serve entire file
		object, err := minioClient.GetObject(ctx, minioClient.MusicBucket, filename)
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusInternalServerError)
			log.Printf("Error getting object: %v", err)
			return
		}
		defer object.Close()

		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, object)
		log.Printf("Served full MinIO file: %s (%d bytes)", filename, fileSize)
		return
	}

	// Parse range header
	ranges, err := parseRange(rangeHeader, fileSize)
	if err != nil || len(ranges) == 0 {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// For simplicity, only handle single range requests
	start := ranges[0].start
	end := ranges[0].end

	// Get object with range
	opts := minio.GetObjectOptions{}
	opts.SetRange(start, end)

	object, err := minioClient.Client.GetObject(ctx, minioClient.MusicBucket, filename, opts)
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		log.Printf("Error getting object with range: %v", err)
		return
	}
	defer object.Close()

	// Set headers for partial content
	contentLength := end - start + 1
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.WriteHeader(http.StatusPartialContent)

	// Copy the requested range
	io.Copy(w, object)
	log.Printf("Served MinIO range: %s (bytes %d-%d/%d)", filename, start, end, fileSize)
}

// ListMinIOMusic returns a list of available music files in MinIO
func ListMinIOMusic(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var files []MediaFile
	objectCh := minioClient.ListObjects(ctx, minioClient.MusicBucket)

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
			Url:         "/gomedia/api/music/" + name,
			ContentType: getContentType(name),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"files": %s}`, formatFileList(files))))
}
