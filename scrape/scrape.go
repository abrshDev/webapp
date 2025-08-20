package scrape

import (
	"context"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

type IGProfileInfo struct {
	Images       []string
	ProfileImage string
	Followers    string
}

// GetIGProfileInfo fetches Instagram profile info safely for container deployment
func GetIGProfileInfo(username string) (*IGProfileInfo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	// Detect Render environment
	if os.Getenv("RENDER") != "" {
		// Connect to Render's headless-shell
		remoteURL := "ws://127.0.0.1:9223/devtools/browser"
		allocatorCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), remoteURL)
		defer allocCancel()
		ctx, cancel = chromedp.NewContext(allocatorCtx)
	} else {
		// Local: launch Chromium manually
		chromePath := os.Getenv("HEADLESS_CHROME_PATH")
		if chromePath == "" {
			chromePath = "/usr/bin/chromium" // default path in Docker
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
			chromedp.Flag("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/116.0.0.0 Safari/537.36"),
		)

		allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancelAlloc()
		ctx, cancel = chromedp.NewContext(allocCtx)
	}

	defer cancel()

	// Increase timeout
	ctx, cancelTimeout := context.WithTimeout(ctx, 120*time.Second)
	defer cancelTimeout()

	var imageurls []string
	var profileImg string
	var followers string

	url := "https://www.instagram.com/" + username
	jscode := `Array.from(document.querySelectorAll('article img')).slice(0,6).map(img => img.src)`

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("header img", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(jscode, &imageurls),
		chromedp.Evaluate(`document.querySelector('header img') ? document.querySelector('header img').src : ''`, &profileImg),
		chromedp.Evaluate(`document.querySelector('header li span[title]') ? document.querySelector('header li span[title]').getAttribute('title') : ''`, &followers),
	)
	if err != nil {
		return nil, err
	}

	return &IGProfileInfo{
		Images:       imageurls,
		ProfileImage: profileImg,
		Followers:    followers,
	}, nil
}
