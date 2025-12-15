# ============================================
# Stage 1: Build the Go application
# ============================================
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Install make command
RUN apk update && apk add make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate swagger docs
RUN make swagger-install
RUN make swagger

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o portfolio-api cmd/api/main.go

# ============================================
# Stage 2: Create minimal runtime image
# ============================================
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs wget

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# Copy the binary from builder
COPY --from=builder /app/portfolio-api .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy swagger docs
COPY --from=builder /app/docs ./docs

# Copy static assets
COPY --from=builder /app/pkg ./pkg

# Create directories with proper permissions
RUN mkdir -p ./data ./pkg/assets/certifications && \
    chown -R appuser:appuser ./data ./pkg

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./portfolio-api"]
