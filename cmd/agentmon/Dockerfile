FROM golang:1.24.1-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o agentmon ./cmd/agentmon

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/agentmon .

CMD ["./agentmon"]
