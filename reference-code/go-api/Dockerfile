# Multi-stage Dockerfile for Creator API
# Stage 1: Build the Go binary
# Stage 2: Create minimal runtime image

# =============================================================================
# Build Stage
# =============================================================================
FROM golang:1.23-alpine AS builder

# Install build dependencies
# - git: for go mod download with private repos
# - ca-certificates: for HTTPS connections
# - tzdata: timezone data for time operations
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy dependency files first (better layer caching)
COPY go.mod go.sum ./

# Download dependencies (cached unless go.mod/go.sum changes)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build binary with optimizations
# CGO_ENABLED=0: static binary (no dynamic dependencies)
# -ldflags="-w -s": strip debug info for smaller binary
# -ldflags="-X main.Version": inject version at build time
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -o creator-api \
    cmd/server/main.go

# =============================================================================
# Runtime Stage
# =============================================================================
FROM alpine:3.19

# Install runtime dependencies
# - ca-certificates: for HTTPS connections
# - tzdata: timezone data
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
# Group: appuser (GID 1000)
# User: appuser (UID 1000)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/creator-api .

# Copy docs directory (needed for Swagger UI if embedded)
COPY --from=builder /build/docs ./docs

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user (security best practice)
USER appuser

# Expose port (documentation only, actual port set via PORT env var)
EXPOSE 34075

# Health check
# Checks /health endpoint every 30s
# Starts checking after 5s (start-period)
# Times out after 3s per check
# Fails after 3 consecutive failures
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-34075}/health || exit 1

# Run the binary
CMD ["./creator-api"]
