package main

import (
	"fmt"
	"log"

	"github.com/oloomoses/go-craper/crawler"
	"github.com/oloomoses/go-craper/parser"
)

func main() {
	url := "https://books.toscrape.com"

	fmt.Println("Fetching books ....")
	data, err := crawler.Fetch(url)

	if err != nil {
		log.Fatal("error fetching books", err)
	}

	fmt.Println("Parsing catalogue ...")
	prods, next, err := parser.ParseCatalogue(data, url)

	if err != nil {
		log.Fatal("Parse catalogue failed", err)
	}

	fmt.Println("Products", len(prods))

	for _, p := range prods[:5] {
		fmt.Printf("%v | %v | %v | %v\n", p.ID, p.Title, p.Price, p.Url)
	}

	if next != "" {
		fmt.Println("Next page", next)
	}
}
