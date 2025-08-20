# Stage 1: Build Go application
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Stage 2: Use lightweight headless-shell image
FROM chromedp/headless-shell:126.0.6478.127
WORKDIR /app
COPY --from=builder /app/app .

# Set environment variables
ENV HEADLESS_CHROME_PATH=/headless-shell/headless-shell
ENV PORT=3000
# Remove hardcoded proxy and IMGBB_KEY for security
# ENV PROXY_URL="http://your-proxy:port"
# ENV IMGBB_KEY="your-key"

EXPOSE 3000
CMD ["./app"]