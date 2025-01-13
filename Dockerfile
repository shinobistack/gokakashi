FROM node:22-alpine3.20 AS frontend

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

ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

COPY --from=frontend /webapp/dist /app/webapp/dist

# Run the tests.CGO is needed for using sqlite3 in tests
RUN CGO_ENABLED=1 go test -v ./...

# Build the Go binary for amd64
RUN CGO_ENABLED=0 GOARCH=amd64 go build -o gokakashi

FROM alpine:3.20

# Set working directory
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/gokakashi /app/gokakashi

# Make sure the binary is executable
RUN chmod +x /app/gokakashi

CMD ["/app/gokakashi"]

# Set the entrypoint to the application binary
ENTRYPOINT ["/app/gokakashi"]
