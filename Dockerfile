# Build stage
FROM golang:1.24.4-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tilokit .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates git nodejs npm

# Create non-root user
RUN addgroup -g 1001 -S tilokit && \
    adduser -S tilokit -u 1001 -G tilokit

# Set working directory
WORKDIR /workspace

# Copy binary from builder stage
COPY --from=builder /app/tilokit /usr/local/bin/tilokit

# Make binary executable
RUN chmod +x /usr/local/bin/tilokit

# Change ownership
RUN chown -R tilokit:tilokit /workspace

# Switch to non-root user
USER tilokit

# Set entrypoint
ENTRYPOINT ["tilokit"]

# Default command
CMD ["--help"]
