# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go workspace files if using go.work (recommended for monorepo)
COPY go.work ./
COPY go.work.sum ./

# Copy module files
COPY worker-service/go.mod worker-service/go.sum ./worker-service/

# Download dependencies
WORKDIR /app/worker-service
RUN go mod download

# Copy full source
WORKDIR /app
COPY . .

# Build worker binary
WORKDIR /app/worker-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker cmd/main.go


# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Install CA certificates (important for HTTPS calls)
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/worker-service/worker .

CMD ["./worker"]
