# Multi-stage build: tiny final image
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Dependencies first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o barkdroid .

# --- Runtime ---
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -h /app barkdroid

USER barkdroid
WORKDIR /app

COPY --from=builder /build/barkdroid .

# Default config location
RUN mkdir -p /app/data

EXPOSE 8080

ENTRYPOINT ["./barkdroid"]
