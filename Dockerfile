# Multi-stage build for smaller image
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and netcat for healthcheck
RUN apk --no-cache add ca-certificates netcat-openbsd

# Create app user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/build/blockchain /app/blockchain

# Create data directory
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose default port
EXPOSE 3000

# Set environment variables
ENV BLOCKCHAIN_DATA_DIR=/app/data
ENV NODE_PORT=3000

# Default command
CMD ["/app/blockchain"]

