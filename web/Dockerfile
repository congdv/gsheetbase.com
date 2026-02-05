# Multi-stage build for Gsheetbase Web Service
# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/ui
COPY web/ui/package*.json ./
RUN npm ci

COPY web/ui/ ./

# Build arg for API URL (will be /api since served from same origin)
ARG VITE_API_BASE_URL=/api
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the web API
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main web/cmd/api/main.go

# Stage 3: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S appuser -u 1001

# Set working directory
WORKDIR /app

# Copy the Go binary
COPY --from=backend-builder /app/main ./main
RUN chmod +x ./main

# Copy the built frontend
COPY --from=frontend-builder /app/ui/dist ./web/ui/dist

# Change ownership to appuser
RUN chown -R appuser:appuser /app

USER appuser

# Expose port (Railway will set PORT env var)
EXPOSE 8080

# Command to run
CMD ["/app/main"]
