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

# Build args for user configuration (defaults match common Linux users)
ARG USER_UID=1000
ARG USER_GID=1000
ARG USERNAME=ct

# Create non-root user and group
RUN addgroup --system --gid ${USER_GID} ${USERNAME} && \
    adduser --system --uid ${USER_UID} --ingroup ${USERNAME} ${USERNAME}

# Create data directory with proper ownership
RUN mkdir -p /app/data && \
    chown -R ${USERNAME}:${USERNAME} /app/data

# Copy binary from builder with proper ownership
COPY --from=builder --chown=${USERNAME}:${USERNAME} /app/server /app/server

# Switch to non-root user
USER ${USERNAME}

# Expose port
EXPOSE 8080

# Run
CMD ["/app/server"]
