package feed

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {

	feed := Print("https://astro-club-ecell-kq371zc1f-shubham-rasal.vercel.app/rss.xml")
	fmt.Print(feed)

}
