# Dockerfile for goKakashi
FROM golang:1.22-alpine

WORKDIR /app

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
