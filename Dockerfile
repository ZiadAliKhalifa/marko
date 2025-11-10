# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app/backend

# Copy go mod files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S -u 1001 -G appuser appuser

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/backend/main .

# Change ownership to non-root user
RUN chown appuser:appuser main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]