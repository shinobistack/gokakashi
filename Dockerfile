# Dockerfile for goKakashi
FROM golang:1.22-bookworm

WORKDIR /app

# Install Docker
RUN  apt-get update \
     apt-get install ca-certificates curl \
     install -m 0755 -d /etc/apt/keyrings \
     curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc \
     chmod a+r /etc/apt/keyrings/docker.asc \
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
       tee /etc/apt/sources.list.d/docker.list > /dev/null \
     apt-get update \
     apt-get install docker-ce-cli


# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o goKakashi ./cmd/goKakashi.go

# Expose ports
EXPOSE 8080
EXPOSE 9090

# Set environment variables
ENV DOCKER_USERNAME="your-dockerhub-username"
ENV DOCKER_PASSWORD="your-dockerhub-password"
ENV LINEAR_API_KEY="your-linear-api-key"

# Run the application
ENTRYPOINT ["./goKakashi"]
