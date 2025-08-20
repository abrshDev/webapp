FROM golang:1.24

# Install dependencies + Chromium
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
    chromium \
    --no-install-recommends && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o app

# Use Chromium for Chromedp
ENV HEADLESS_CHROME_PATH=/usr/bin/chromium

ENV PORT=3000
ENV PROXY_URL="http://192.168.120.122:8080"
ENV IMGBB_KEY="904775b3a745b64f07d3f6dff7407701"

EXPOSE 3000
CMD ["./app"]
