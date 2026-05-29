# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

# Install CA certs for module downloads
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /user-service ./main.go

# Runtime stage
FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /user-service /user-service

ENV PORT=8080

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/user-service"]
