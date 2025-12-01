package main

import (
	"fmt"
	"log"

	"github.com/oloomoses/go-craper/crawler"
)

func main() {
	url := "https://books.toscrape.com/catalogue/category/books/travel_2/index.html"

	fmt.Println("Fetching books ....")
	data, err := crawler.Fetch(url)

	if err != nil {
		log.Fatal("error fetching books", err)
	}

	fmt.Println(data)

}
