# Multi-stage build for Gsheetbase Worker Service
# Stage 1: Build the Go worker
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the worker service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main worker/cmd/api/main.go

# Stage 2: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S appuser -u 1001

# Set working directory
WORKDIR /app

# Copy the Go binary
COPY --from=builder /app/main ./main
RUN chmod +x ./main

# Change ownership to appuser
RUN chown -R appuser:appuser /app

USER appuser

# Expose port (Railway will set PORT env var, but worker typically doesn't serve HTTP)
EXPOSE 8080

# Command to run
CMD ["/app/main"]
