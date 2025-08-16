package scrape

import (
	"context"
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
	// Replace this with your proxy (mobile or residential)
	proxy := "http://192.168.120.122:8080"

	// Create Chromedp ExecAllocator with safe flags + proxy
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.ProxyServer(proxy),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	// Timeout to prevent hanging
	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
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
