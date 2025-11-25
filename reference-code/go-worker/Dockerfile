# Build stage
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

# GitHub token for private repos
ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/orchard9

# Copy go mod files
COPY go.mod go.sum ./

# Configure git to use token and download modules
RUN if [ -n "$GITHUB_TOKEN" ]; then \
    git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"; \
    fi && \
    go mod download && \
    git config --global --unset url."https://${GITHUB_TOKEN}@github.com/".insteadOf || true

# Copy source code
COPY . .

# Build binary with service name for squiddy convention
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG SERVICE_NAME=email-worker
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -installsuffix cgo -o ${SERVICE_NAME} ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Use /app as base directory (squiddy convention)
WORKDIR /app

# Create bin directory
RUN mkdir -p /app/bin

# Copy binary to expected location with service name
ARG SERVICE_NAME=email-worker
COPY --from=builder /build/${SERVICE_NAME} /app/bin/${SERVICE_NAME}

# Make executable
RUN chmod +x /app/bin/${SERVICE_NAME}

EXPOSE 8080

# Use full path in CMD
CMD ["/app/bin/email-worker"]
