# Multi-stage build for Go backend
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/seed .

# Copy swagger docs
COPY --from=builder /app/docs ./docs

# Create directories
RUN mkdir -p /app/uploads /app/logs

# Expose port
EXPOSE 8080

# Run
CMD ["./main"]
