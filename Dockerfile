FROM golang:1.25.0-alpine

WORKDIR /usr/src/app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build binary as root user
# Create data dir and non-root user
RUN go build -o /usr/local/bin/whoknows_variations ./cmd && \
    addgroup --system appgroup && \
    adduser --system appuser --ingroup appgroup && \
    mkdir -p /usr/src/app/data && \
    chown -R appuser:appgroup /usr/src/app

# Switch to non-root user
USER appuser

EXPOSE 8080
CMD ["whoknows_variations"]
