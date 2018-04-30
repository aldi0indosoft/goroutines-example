package crawler

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Crawl is Crawl
func Crawl(uri string, chURL chan string, chDone chan bool) {
	resp, err := http.Get(uri)
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
			ok, uri := GetHref(t)
			if !ok {
				continue
			}
			if hasProto := strings.Index(uri, "http") == 0; hasProto {
				chURL <- uri
			}
		}
	}
}

// CrawlWithProxy is CrawlWithProxy
func CrawlWithProxy(uri string, chURL chan string, chDone chan bool) {
	// create client with proxy
	u, err := url.Parse("http://proxy.indo-soft.com:443")
	if err != nil {
		panic(err)
	}
	tr := &http.Transport{
		Proxy: http.ProxyURL(u),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	client := &http.Client{Transport: tr}

	// begin crawl
	resp, err := client.Get(uri)
	defer func() {
		chDone <- true
	}()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	b := resp.Body
	defer b.Close()

	// fix html first
	root, err := html.Parse(b)

	if err != nil {
		log.Fatal(err)
	}

	var cbuff bytes.Buffer
	html.Render(&cbuff, root)
	fixed := cbuff.String()
	// fmt.Println("response: ", fixed)

	z := html.NewTokenizer(bytes.NewBufferString(fixed))

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
			ok, uri := GetHref(t)
			if !ok {
				continue
			}
			chURL <- uri
			if hasProto := strings.Index(uri, "http") == 0; hasProto {
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
