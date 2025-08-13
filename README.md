# Insta Scraper Go + Fiber + Chromedp + Replit + ngrok

## How to run on Replit

1. Add your ngrok authtoken to `ngrok.yml` (replace `<YOUR_NGROK_AUTHTOKEN>`).
2. Replit will run `go run main.go` automatically.
3. In the Shell, run:

   ngrok start --all --config=ngrok.yml

4. Use the public ngrok URL to access your Fiber server endpoints.

---

- The server runs on port 3000.
- Example endpoint: `/images/<username>`

# Directory Structure

└── scrape
    ├── scrape.go
├── .replit
├── go.mod
├── go.sum
├── main.go
├── ngrok.yml
├── README.md
├── replit.nix

# End Directory Structure