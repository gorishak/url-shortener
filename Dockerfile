# Build stage
FROM golang:1.25.1-alpine AS builder

# Install build dependencies for CGO (required for SQLite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled (required for SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -o url-shortener ./cmd/url-shortener

# Runtime stage
FROM alpine:latest

# Install SQLite dependencies
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/url-shortener .

# Copy config file
COPY config/prod.yaml /app/config/prod.yaml

# Create storage directory
RUN mkdir -p /app/storage

# Expose port
EXPOSE 8082

# Set environment variable
ENV CONFIG_PATH=/app/config/prod.yaml

# Run the application
CMD ["./url-shortener"]

