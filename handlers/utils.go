package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// MediaFile represents a media file with metadata
type MediaFile struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Path        string `json:"path,omitempty"` // Omit path in JSON response for security if desired, but keeping for now as per previous
	Url         string `json:"url"`
	ContentType string `json:"contentType"`
}

// listMediaFiles lists all files in a directory and generates metadata
func listMediaFiles(dir string, baseUrl string, getContentType func(string) string) ([]MediaFile, error) {
	var files []MediaFile

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := entry.Name()
		files = append(files, MediaFile{
			Name:        name,
			Size:        info.Size(),
			Path:        filepath.Join(dir, name),
			Url:         baseUrl + "/" + name,
			ContentType: getContentType(name),
		})
	}

	return files, nil
}

// formatFileList formats a list of files as JSON
func formatFileList(files []MediaFile) string {
	if len(files) == 0 {
		return "[]"
	}

	data, err := json.Marshal(files)
	if err != nil {
		return "[]"
	}

	return string(data)
}

// getContentType returns the MIME type based on file extension
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	case ".flac":
		return "audio/flac"
	default:
		return "application/octet-stream"
	}
}

// getImageContentType returns the MIME type for images
func getImageContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".bmp":
		return "image/bmp"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

// byteRange represents a byte range
type byteRange struct {
	start int64
	end   int64
}

// parseRange parses HTTP Range header
func parseRange(rangeHeader string, fileSize int64) ([]byteRange, error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return nil, fmt.Errorf("invalid range header")
	}

	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeSpec, "-")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range format")
	}

	var start, end int64
	var err error

	if parts[0] == "" {
		// Suffix range: -500 (last 500 bytes)
		end = fileSize - 1
		start = fileSize - parseInt64(parts[1], 0)
		if start < 0 {
			start = 0
		}
	} else if parts[1] == "" {
		// Start only: 500- (from byte 500 to end)
		start = parseInt64(parts[0], 0)
		end = fileSize - 1
	} else {
		// Both start and end: 500-999
		start = parseInt64(parts[0], 0)
		end = parseInt64(parts[1], fileSize-1)
	}

	if start > end || start < 0 || end >= fileSize {
		return nil, fmt.Errorf("invalid range values")
	}

	return []byteRange{{start: start, end: end}}, err
}

// parseInt64 safely parses a string to int64 with a default value
func parseInt64(s string, defaultVal int64) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultVal
	}
	return val
}
