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

func GetIGProfileInfo(username string) (*IGProfileInfo, error) {
	// Replace with your phone's private IP and port from Every Proxy
	proxy := "http://192.168.120.122:8080"

	// Create Chrome options to use the proxy
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(proxy), // Set proxy
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

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
