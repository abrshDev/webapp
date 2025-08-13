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
	ctx, cancel := chromedp.NewContext(context.Background())
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
