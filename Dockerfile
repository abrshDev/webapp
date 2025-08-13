# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

# If you have static files, copy them as well:
# COPY --from=builder /app/static ./static

ENV PORT=8080