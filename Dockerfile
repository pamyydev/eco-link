# Multi-stage Dockerfile: build Go binary and serve pre-built frontend dist

# Builder stage for Go
FROM golang:1.25-bullseye as builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/eco-link ./cmd/webapp

# Final image: small Debian slim with CA certs
FROM debian:stable-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
# Create app user
RUN useradd -m appuser || true

WORKDIR /app
COPY --from=builder /app/eco-link /usr/local/bin/eco-link
# Copy frontend dist if present
COPY frontend/dist /app/frontend/dist

ENV PORT=8081
EXPOSE 8081

USER appuser
ENTRYPOINT ["/usr/local/bin/eco-link"]
