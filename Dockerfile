FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    zstd \
    tar

# Set working directory
WORKDIR /app

# Copy source code
COPY src/ ./src/
COPY go.mod ./
COPY go.sum ./

# Build the Go application
RUN cd src && go build -o ../main main.go

# Final image
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache \
    git \
    zstd \
    tar \
    && adduser -D -g '' appuser

# Switch to non-root user
USER appuser
WORKDIR /home/appuser

# Copy built binary
COPY --from=builder /app/main /app/main

# Set entrypoint
ENTRYPOINT ["/app/main"]