package crawler

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Crawl is Crawl
func Crawl(url string, chURL chan string, chDone chan bool) {
	resp, err := http.Get(url)
	defer func() {
		chDone <- true
	}()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	b := resp.Body
	defer b.Close()

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t := z.Token()
			if isAnchor := t.Data == "a"; !isAnchor {
				continue
			}
			ok, url := GetHref(t)
			if !ok {
				continue
			}
			if hasProto := strings.Index(url, "http") == 0; hasProto {
				chURL <- url
			}
		}
	}
}

// GetHref is GetHref
func GetHref(t html.Token) (bool, string) {
	for _, attr := range t.Attr {
		if attr.Key == "href" {
			return true, attr.Val
		}
	}
	return false, ""
}
