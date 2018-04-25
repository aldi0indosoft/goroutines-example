package main

import (
	"crawler"
	"fmt"
)

func main() {
	urls := []string{
		"https://www.olx.co.id/all-results/q-iphone-lock/?search%5Border%5D=filter_float_price%3Aasc",
		"https://www.bukalapak.com/products?utf8=%E2%9C%93&source=navbar&from=omnisearch&search_source=omnisearch_organic&search%5Bkeywords%5D=iphone+lock",
	}
	links := make(map[string]bool)

	chURL := make(chan string)
	chDone := make(chan bool)

	for _, url := range urls {
		fmt.Println(url)
		go crawler.Crawl(url, chURL, chDone)
	}
	for c := 0; c < len(urls); {
		select {
		case url := <-chURL:
			links[url] = true
		case <-chDone:
			c++
		}
	}

	for link := range links {
		fmt.Println("Link: ", link)
	}
	close(chURL)
}
