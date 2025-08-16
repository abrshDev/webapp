# Use Chromedp headless Chrome + Go
FROM golang:1.24

# Install dependencies for Chromedp
RUN apt-get update && apt-get install -y \
    ca-certificates \
    fonts-liberation \
    libappindicator3-1 \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcups2 \
    libdbus-1-3 \
    libdrm2 \
    libgbm1 \
    libnspr4 \
    libnss3 \
    libx11-xcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    xdg-utils \
    --no-install-recommends && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy Go modules and download
COPY go.mod go.sum ./
RUN go mod download

# Copy all files
COPY . .

# Build the Go app
RUN go build -o app

# Set environment
ENV PORT=3000

# Expose port
EXPOSE 3000

# Run the app
CMD ["./app"]
