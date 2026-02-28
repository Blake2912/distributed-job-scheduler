# -------------------------
# Build stage
# -------------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Copy go workspace files from root
COPY go.work go.work.sum ./

# Copy only go.mod files first (for layer caching)
COPY worker-service/go.mod worker-service/go.sum ./worker-service/
COPY common/go.mod common/go.sum ./common/

# Download dependencies
WORKDIR /app/worker-service
RUN go mod download

# Now copy full source
WORKDIR /app
COPY . .

# Build worker binary
WORKDIR /app/worker-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker cmd/main.go


# -------------------------
# Runtime stage
# -------------------------
FROM alpine:latest

WORKDIR /root/

# Install CA certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy built binary
COPY --from=builder /app/worker-service/worker .

# Run binary
CMD ["./worker"]
