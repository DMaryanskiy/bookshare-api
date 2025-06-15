# Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API binary
RUN go build -o /bookshare-api ./cmd/api

# Build Worker binary
RUN go build -o /bookshare-worker ./cmd/worker

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /bookshare-api /bookshare-api
COPY --from=builder /bookshare-worker /bookshare-worker

EXPOSE 8080

ENTRYPOINT ["/bookshare-api"]
