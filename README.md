# 🎵 Go Media Streaming Backend with MinIO

A high-performance Go backend for streaming music and images exclusively from MinIO object storage, featuring HTTP range request support, intelligent caching, and CORS handling.

## ✨ Features

### Capabilities
- **MinIO Native**: Stream directly from S3-compatible object storage
- **Music Streaming**: HTTP range request support for seeking and partial content delivery
- **Image Streaming**: Efficient caching using ETags and Last-Modified headers
- **REST API**: Clean `/api/` endpoints for listing and streaming
- **CORS Enabled**: Cross-origin resource sharing for web clients
- **Request Logging**: Comprehensive HTTP request logging

##  Quick Start

### Prerequisites

- Go 1.21 or higher
- MinIO server (local or remote)

### Installation

1. Navigate to the project directory:
```bash
cd /Users/duylong/Code/PersonalProject/MediaBackend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure environment variables (copy from example):
```bash
cp .env.example .env
# Edit .env with your MinIO credentials
```

4. Run the server:
```bash
go run main.go
```
*Or simply `go run .`*

The server will start on `http://localhost:8080`.

## ⚙️ Configuration

Create a `.env` file or set environment variables:

```bash
# Server Configuration
PORT=8080

# MinIO Configuration
MINIO_ENDPOINT=localhost:9098
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_MUSIC_BUCKET=music
MINIO_IMAGE_BUCKET=images
```

## 📡 API Endpoints

### Music

- **List Music**: `GET /api/music`
  - Returns JSON list of tracks with metadata.

- **Stream Music**: `GET /api/music/{filename}`
  - Supports HTTP range requests for seeking.
  - Example: `http://localhost:8022/api/music/song.mp3`

### Images

- **List Images**: `GET /api/images`
  - Returns JSON list of images with metadata.

- **Stream Image**: `GET /api/images/{filename}`
  - Supports caching with ETags.
  - Example: `http://localhost:8080/api/images/photo.jpg`

### Utility

- **Health Check**: `GET /health`
- **Test Client**: `GET /`

## 📁 Project Structure

```
MediaBackend/
├── main.go                 # Server entry point
├── handlers/
│   ├── minio_music.go     # MinIO music streaming
│   ├── minio_image.go     # MinIO image streaming
│   ├── utils.go           # Utility functions & Structs
│   └── client.go          # Test client HTML
├── middleware/
│   ├── cors.go            # CORS middleware
│   └── logging.go         # Request logging
├── minio/
│   └── config.go          # MinIO client configuration
├── go.mod
├── .env.example
└── README.md
```

## 🔐 Security Features

- **Path Traversal Protection**: Prevents access outside specific buckets
- **CORS Configuration**: Configurable cross-origin access
- **Input Sanitization**: Validates and sanitizes all file paths
- **MinIO Authentication**: Secure credential-based access

## 🐳 Running with MinIO

### Start MinIO with Docker:

```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  minio/minio server /data --console-address ":9001"
```

### Upload files to MinIO:

```bash
# Using MinIO Client (mc)
mc alias set local http://localhost:9000 minioadmin minioadmin
mc mb local/music
mc mb local/images
mc cp song.mp3 local/music/
mc cp photo.jpg local/images/
```

##  Building for Production

```bash
# Build binary
go build -o mediaserver

# Run binary
./mediaserver
```
