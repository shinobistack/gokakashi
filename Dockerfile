FROM node:22-alpine3.20 AS frontend

# Set the working directory
WORKDIR /webapp

COPY webapp/package.json webapp/package-lock.json ./

RUN npm install

COPY webapp/ /webapp/

RUN npm run build

# Stage 1: Build the Go binary using Alpine
FROM golang:1.23-alpine AS builder

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

# Set working directory
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/gokakashi /app/gokakashi

# Make sure the binary is executable
RUN chmod +x /app/gokakashi

# Set the entrypoint to the application binary
ENTRYPOINT ["/app/gokakashi"]
