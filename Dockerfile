# Build stage
FROM golang:1.25.0-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/auth-service ./main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/bin/auth-service .

# Copy the .env file (make sure to include it in your build context and .dockerignore if necessary)
COPY .env /app/bin/.env

# Expose port (default is 8080)
EXPOSE 8080

# Run the application
CMD ["./auth-service"]
