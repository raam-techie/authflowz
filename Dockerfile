# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /server

# Copy go mod and sum files to leverage Docker cache
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if go.mod and go.sum are unchanged
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o /server/app ./main.go

# Stage 2: Final stage for production
FROM alpine:latest AS production

# Install tzdata package for time zone configuration
RUN apk --no-cache add tzdata

# Set environment variable for the time zone, this will be used at runtime
ENV TZ=Asia/Kolkata

RUN echo $TZ

# Configure the system's time zone using the TZ variable
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Install curl for health checks. This is required.
RUN apk --no-cache add curl

# Set the working directory inside the container
WORKDIR /server

# Copy the built executable from the builder stage
COPY --from=builder /server/app .

# Copy environment variables file to /server
COPY .env /server/.env

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the built application
CMD ["./app"]

