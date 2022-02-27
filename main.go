package main

import (
	"log"
	"net/url"
	"os"

	"github.com/gocolly/colly/v2"
)

func main() {
	root, err := url.Parse("https://geohot.github.io/blog/")
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector(
		colly.UserAgent("googlebot"),
		colly.CacheDir("./.cache"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Printf("visiting %s\n", r.URL)
	})
	// Extract text off of webpage
	c.OnResponse(func(r *colly.Response) {
		os.WriteFile("./output", r.Body, 0644)
	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(link)
	})

	c.Visit(root.String())
}
