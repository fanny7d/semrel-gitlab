FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG VERSION=DEV
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X main.version=${VERSION}" -o release .

# Final stage
FROM alpine:3.21.3

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/release .

# Run the binary
ENTRYPOINT ["./release"] 