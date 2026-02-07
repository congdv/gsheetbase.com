# Multi-stage build for Gsheetbase Web Service
# Stage 1: Build the frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/ui
COPY web/ui/package*.json ./
RUN npm ci

COPY web/ui/ ./

# Build args for Vite environment variables (must be provided at build time)
ARG VITE_API_BASE_URL
ARG VITE_WORKER_BASE_URL
ARG VITE_LANDING_PAGE_URL
ARG VITE_FORCE_PROD

# Set environment variables for the build
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
ENV VITE_WORKER_BASE_URL=$VITE_WORKER_BASE_URL
ENV VITE_LANDING_PAGE_URL=$VITE_LANDING_PAGE_URL
ENV VITE_FORCE_PROD=$VITE_FORCE_PROD

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
