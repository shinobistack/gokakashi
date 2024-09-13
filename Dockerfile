# Stage 1: Build the Go binary
FROM golang:1.22-bookworm AS builder

# Ensure the build fails on any command failure
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

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

# Stage 2: Final image for running the application
FROM debian:bookworm-slim

# Ensure the build fails on any command failure
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Set working directory
WORKDIR /app

# Install dependencies: Docker CLI and Trivy
RUN apt-get update && \
    apt-get install -y ca-certificates curl gnupg lsb-release && \
    install -m 0755 -d /etc/apt/keyrings && \
    curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc && \
    chmod a+r /etc/apt/keyrings/docker.asc && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
    bullseye stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null && \
    apt-get update && \
    apt-get install -y docker-ce-cli
# Install Trivy v0.55.0
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/install.sh | sh -s -- -b /usr/local/bin

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
