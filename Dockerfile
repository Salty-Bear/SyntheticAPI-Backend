# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the Go binary
RUN go build -o pubsub main.go

# Run stage
FROM alpine:3.18
WORKDIR /app

# Copy the built binary
COPY --from=builder /app/pubsub .

# Expose port
EXPOSE 3000

# Run the binary
CMD ["./pubsub"]