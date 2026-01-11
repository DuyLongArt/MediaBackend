package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"MediaBackend/handlers"
	"MediaBackend/middleware"
	minioClient "MediaBackend/minio"
)

func main() {
	// Initialize MinIO client
	if err := minioClient.InitMinIO(); err != nil {
		log.Printf("⚠️  MinIO initialization failed: %v", err)
		log.Printf("⚠️  MinIO endpoints will not be available")
	}

	// Set up routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// MinIO streaming endpoints (REST API)
	mux.HandleFunc("/api/music/", handlers.StreamMinIOMusic)
	mux.HandleFunc("/api/images/", handlers.StreamMinIOImage)
	mux.HandleFunc("/api/music", handlers.ListMinIOMusic)
	mux.HandleFunc("/api/images", handlers.ListMinIOImages)

	// Serve test client
	mux.HandleFunc("/", handlers.ServeTestClient)

	// Apply middleware
	handler := middleware.CORS(middleware.Logging(mux))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8022"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🎵 Media Streaming Server starting on http://localhost%s", addr)
	log.Printf("☁️  MinIO Music: %s", minioClient.MusicBucket)
	log.Printf("☁️  MinIO Images: %s", minioClient.ImageBucket)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
