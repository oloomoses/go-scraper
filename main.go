package main

import (
	"github.com/oloomoses/go-craper/crawler"
)

func main() {
	url := "https://books.toscrape.com"

	crwl := crawler.New(url)
	crwl.Crawl(url, 100)

}
