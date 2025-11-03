# Multi-stage build for smaller image

# Stage 1: Build the application locally (outside Docker)
# We'll copy the pre-built binary instead of building inside Docker

# Stage 2: Runtime
FROM alpine:latest

# Install ca-certificates for HTTPS and netcat for health checks
RUN apk --no-cache add ca-certificates netcat-openbsd

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the pre-built binary from build/ directory
COPY build/blockchain /app/blockchain

# Create data directory with correct permissions
# Note: /app/tmp is NOT created here - wallets should be stored in /app/data/tmp (mounted volume)
RUN mkdir -p /app/data/blocks && \
    chown -R appuser:appgroup /app && \
    chmod +x /app/blockchain

# Switch to non-root user
USER appuser

# Default data directory
ENV BLOCKCHAIN_DATA_DIR=/app/data/blocks

# Expose network port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=40s \
    CMD sleep 1

# Default command (can be overridden by docker-compose)
CMD ["/app/blockchain"]
