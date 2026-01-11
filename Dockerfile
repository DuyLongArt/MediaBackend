# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# -o mediaserver: output binary name
# -ldflags="-s -w": strip debug information for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mediaserver .

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/mediaserver .
COPY --from=builder /app/.env.example .env.example

# Expose port
EXPOSE 8022

# Environment variables with defaults (can be overridden)
ENV PORT=8022
ENV MINIO_ENDPOINT=minio:9000
ENV MINIO_USE_SSL=false

# Run the binary
CMD ["./mediaserver"]
