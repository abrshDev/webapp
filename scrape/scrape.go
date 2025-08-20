package scrape

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

type IGProfileInfo struct {
	Images       []string
	ProfileImage string
	Followers    string
}

func GetIGProfileInfo(username string) (*IGProfileInfo, error) {
	chromePath := os.Getenv("HEADLESS_CHROME_PATH")
	if chromePath == "" {
		chromePath = "/usr/bin/chromium"
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("no-zygote", true),
		chromedp.Flag("single-process", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/116.0.0.0 Safari/537.36"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	ctx, cancelTimeout := context.WithTimeout(ctx, 180*time.Second)
	defer cancelTimeout()

	var imageurls []string
	var profileImg string
	var followers string

	url := "https://www.instagram.com/" + username

	retries := 3
	for i := 0; i < retries; i++ {
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(7*time.Second), // longer wait for JS content
			chromedp.Evaluate(`Array.from(document.querySelectorAll('article img')).slice(0,6).map(img => img.src)`, &imageurls),
			chromedp.Evaluate(`document.querySelector('header img') ? document.querySelector('header img').src : ''`, &profileImg),
			chromedp.Evaluate(`document.querySelector('header li span[title]') ? document.querySelector('header li span[title]').getAttribute('title') : ''`, &followers),
		)
		if err == nil && len(imageurls) > 0 {
			break
		}
		log.Printf("Attempt %d failed, retrying...", i+1)
		time.Sleep(3 * time.Second)
	}

	// If still empty, log full HTML for debugging
	if len(imageurls) == 0 && profileImg == "" && followers == "" {
		var html string
		chromedp.Run(ctx, chromedp.OuterHTML("html", &html, chromedp.ByQuery))
		log.Println("Scraping failed, full page HTML:", html)
	}

	return &IGProfileInfo{
		Images:       imageurls,
		ProfileImage: profileImg,
		Followers:    followers,
	}, nil
}
