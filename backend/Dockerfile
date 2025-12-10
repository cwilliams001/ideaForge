# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies (needed for sqlite3)
RUN apk add --no-cache gcc musl-dev

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO enabled for sqlite
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Run
CMD ["/app/server"]
