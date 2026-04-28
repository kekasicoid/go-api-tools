# STAGE 1: Builder
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git build-base

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app as a static binary with size optimization
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o main ./cmd/server/main.go

# STAGE 2: Final Image
FROM alpine:3.19

# Add CA certificates for HTTPS requests and timezone data
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the pre-built binary
COPY --from=builder /app/main .

# Create a non-root user for security
RUN adduser -D appuser && \
    chown -R appuser:appuser /app

USER appuser

# Expose port (default for the app)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
