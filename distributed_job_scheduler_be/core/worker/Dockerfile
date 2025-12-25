# -------- Build stage --------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY worker/go.mod worker/go.sum ./
RUN go mod download

COPY worker/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker

# -------- Runtime stage --------
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/worker .

USER nonroot:nonroot
CMD ["./worker"]
