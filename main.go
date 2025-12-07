package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oloomoses/go-craper/crawler"
)

func main() {
	url := "https://books.toscrape.com/catalogue/page-5.html"

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	crwl := crawler.New(url)
	// crwl.Crawl(url, 100)

	start := time.Now()
	prods, err := crwl.CrawlConcurrently(ctx, url, 60, 35)
	// urls, _ := crwl.DiscoverPages(url, 10)

	if err != nil {
		log.Panic("Failed to fetch products")
	}

	fmt.Printf("\nDone! %d books in %.2f seconds\n", len(prods), time.Since(start).Seconds())

	// fmt.Println(prods)

	// fmt.Println(urls)

}
