# Stage 1: Build the Go binary using Alpine
FROM golang:1.22-alpine AS builder

# Ensure the build fails on any command failure
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# Install build dependencies
RUN apk add --no-cache git bash

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source code
COPY . .

# Build the Go binary for amd64
RUN GOARCH=amd64 go build -o goKakashi ./cmd/goKakashi.go

# Stage 2: Final image for running the application with Alpine
FROM alpine:3.18

# Ensure the build fails on any command failure
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# Install Docker CLI and other dependencies
RUN apk add --no-cache docker-cli curl bash ca-certificates

# Install Trivy
RUN curl -sfL https://github.com/aquasecurity/trivy/releases/download/v0.18.3/trivy_0.18.3_Linux-64bit.tar.gz | tar -xz -C /usr/local/bin

# Set working directory
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/goKakashi /app/goKakashi

# Expose ports
EXPOSE 8080
EXPOSE 9090

# Set environment variables
ENV DOCKER_USERNAME="your-dockerhub-username"
ENV DOCKER_PASSWORD="your-dockerhub-password"
ENV LINEAR_API_KEY="your-linear-api-key"

# Make sure the binary is executable
RUN chmod +x /app/goKakashi

# Set the entrypoint to the application binary
ENTRYPOINT ["/app/goKakashi"]
