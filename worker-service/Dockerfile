# -------------------------
# Build stage
# -------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy only worker module files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy remaining source
COPY . .

# Build worker binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker cmd/main.go


# -------------------------
# Runtime stage
# -------------------------
FROM alpine:latest

WORKDIR /root/
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/worker .

CMD ["./worker"]
