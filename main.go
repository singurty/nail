package main

import (
	"log"
	"net/url"

	"github.com/gocolly/colly/v2"
)

func main() {
	root, err := url.Parse("https://geohot.github.io/blog/")
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector(
		colly.CacheDir("./.cache"),
		colly.AllowedDomains(root.Host),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Printf("scraping %s\n", link)
	})

	c.Visit(root.String())
}
