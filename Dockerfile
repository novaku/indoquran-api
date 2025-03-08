# Stage 1: Build the Go binary
FROM golang:alpine AS builder

# set the environment variable
ENV ENV=docker

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache build-base git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download vendor dependencies
RUN go mod vendor

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o bin/indoquran cmd/api/main.go


# Stage 2: Create a small image with the Go binary
FROM alpine:latest

ENV ENV=docker

# Set working directory in the new container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/indoquran .
COPY --from=builder /app/internal/config/. /app/internal/config/

# # Expose the application on port 8080 (adjust if needed)
EXPOSE 8090

# # Command to run the Go binary
CMD ["./indoquran"]
