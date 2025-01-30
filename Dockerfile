FROM node:22-alpine3.20 AS frontend
WORKDIR /webapp
COPY webapp/package.json webapp/package-lock.json ./
RUN npm install
COPY webapp/ /webapp/
RUN npm run build

FROM golang:1.23-alpine AS builder
SHELL ["/bin/ash", "-o", "pipefail", "-c"]
RUN apk add --no-cache git bash gcc sqlite-dev musl-dev libc-dev
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /webapp/dist /app/webapp/dist
# CGO is needed for using sqlite3 in tests
RUN CGO_ENABLED=1 go test -v ./...
RUN CGO_ENABLED=0 GOARCH=amd64 go build -o gokakashi

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/gokakashi /app/gokakashi
RUN apk add --no-cache libc6-compat libseccomp gcompat
RUN wget -qO- https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh |sh -s -- -b /usr/local/bin v0.58.1
RUN chmod +x /app/gokakashi
CMD ["/app/gokakashi"]
ENTRYPOINT ["/app/gokakashi"]
