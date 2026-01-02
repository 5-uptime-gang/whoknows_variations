# ---- Build stage ----
FROM golang:1.25.0-alpine AS builder
RUN apk add --no-cache git
LABEL org.opencontainers.image.source="https://example.com/whoknows_variations" \
      org.opencontainers.image.description="Build stage for whoknows variations service"
WORKDIR /usr/src/app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOMAXPROCS=1 \
    GOMEMLIMIT=200MiB

RUN go build -o /usr/local/bin/whoknows_variations ./cmd

# ---- Runtime stage ----
FROM alpine:3.20
WORKDIR /usr/src/app

# Create data dir (as root)
RUN mkdir -p /usr/src/app/data

# Copy binary and public assets
COPY --from=builder /usr/local/bin/whoknows_variations /usr/local/bin/whoknows_variations
COPY --from=builder /usr/src/app/public /usr/src/app/public

EXPOSE 8080

CMD ["whoknows_variations"]
