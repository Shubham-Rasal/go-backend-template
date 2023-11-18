package feed

import (
	"github.com/mmcdole/gofeed"
)

func Print(url string) *gofeed.Feed {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	return feed
}
