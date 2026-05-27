# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# Install git for go module downloads
RUN apk add --no-cache git

WORKDIR /app

# Cache dependency downloads separately from source changes
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o assignment-service ./cmd/main.go

# ── Stage 2: Final ────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/assignment-service .

# Copy .env if present (optional; real config via env vars)
COPY .env* ./

EXPOSE 8083

ENTRYPOINT ["./assignment-service"]
