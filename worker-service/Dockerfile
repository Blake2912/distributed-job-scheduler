# -------------------------
# Build stage
# -------------------------
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy workspace file
COPY go.work go.work.sum ./

# Copy entire repository (needed for workspace resolution)
COPY . .

# Build worker
WORKDIR /app/worker-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker cmd/main.go


# -------------------------
# Runtime stage
# -------------------------
FROM alpine:latest

WORKDIR /root/
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/worker-service/worker .

CMD ["./worker"]
