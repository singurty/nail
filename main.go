package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector(
		colly.CacheDir("./.cache"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("scraping %s\n", link)
		e.Request.Visit(link)
	})

	c.Visit("https://geohot.github.io/blog/")
}
