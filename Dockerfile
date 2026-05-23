# Multi-stage build: compile Go binary, then copy to minimal image
FROM golang:1.22-alpine AS builder
WORKDIR /app

# cache modules
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org && go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-s -w' -o /server ./cmd/server

FROM alpine:3.18
RUN adduser -D -H appuser
COPY --from=builder /server /server
USER appuser
EXPOSE 8080
ENTRYPOINT ["/server"]
