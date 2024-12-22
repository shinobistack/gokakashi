FROM node:18-alpine3.20 AS frontend

# Set the working directory
WORKDIR /webapp

COPY webapp/package.json webapp/package-lock.json ./

RUN npm install

COPY webapp/ /webapp/

RUN npm run build

# Stage 1: Build the Go binary using Alpine
FROM golang:1.23-alpine AS builder

# Ensure the build fails on any command failure
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# Install build dependencies
RUN apk add --no-cache git bash gcc sqlite-dev musl-dev libc-dev

# Set CGO_ENABLED for sqlite3 compatibility
# ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source code
COPY . .

COPY --from=frontend /webapp/dist /app/webapp/dist

# Run the tests
RUN go test -v ./...

# Build the Go binary for amd64
RUN GOARCH=amd64 go build -o gokakashi

FROM alpine:3.20

# Ensure the build fails on any command failure
SHELL ["/bin/ash", "-o", "pipefail", "-c"]

# Install Docker CLI and other dependencies
RUN apk add --no-cache docker-cli curl bash ca-certificates python3

# Install Trivy
RUN curl -sfL https://github.com/aquasecurity/trivy/releases/download/v0.55.1/trivy_0.55.1_Linux-64bit.tar.gz | tar -xz -C /usr/local/bin

# Install minimal gcloud CLI
RUN curl -sS -o /tmp/google-cloud-cli.tar.gz https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-494.0.0-linux-x86_64.tar.gz \
    && tar -xvf /tmp/google-cloud-cli.tar.gz \
    && ./google-cloud-sdk/install.sh --quiet

# Add gcloud to the PATH
ENV PATH=$PATH:/google-cloud-sdk/bin

# Set working directory
WORKDIR /app

RUN mkdir -p /app/website

# Copy the Go binary from the builder stage
COPY --from=builder /app/gokakashi /app/gokakashi

# Expose ports
EXPOSE 8080
EXPOSE 9090

# Set environment variables
ENV DOCKER_USERNAME="your-dockerhub-username"
ENV DOCKER_PASSWORD="your-dockerhub-password"
ENV LINEAR_API_KEY="your-linear-api-key"

# Make sure the binary is executable
RUN chmod +x /app/gokakashi

# Set the entrypoint to the application binary
ENTRYPOINT ["/app/gokakashi"]
