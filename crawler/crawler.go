package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/oloomoses/go-craper/model"
	"github.com/oloomoses/go-craper/parser"
	"github.com/oloomoses/go-craper/util"
)

type Crawl struct {
	BaseUrl string
	Visited map[string]bool
}

func New(baseUrl string) *Crawl {
	return &Crawl{
		BaseUrl: baseUrl,
		Visited: make(map[string]bool),
	}
}

func (c *Crawl) Crawl(startUrl string, maxPages int) ([]model.Product, error) {

	var allProducts []model.Product

	currentUrl := startUrl

	pages := 0

	for currentUrl != "" && pages < maxPages {
		if c.Visited[currentUrl] {
			break
		}

		c.Visited[currentUrl] = true

		pages += 1

		fmt.Printf("[%d] Crawling: %s\n", pages, currentUrl)

		products, nextHref, err := c.fetchAndParse(currentUrl)

		if err != nil {
			log.Fatal("Fetching Error, ", err)
			break
		}

		allProducts = append(allProducts, products...)

		fmt.Printf("Foud %d products (total: %d) \n", len(products), len(allProducts))

		if nextHref == "" {
			break
		}

		currentUrl, err = util.Resolve(currentUrl, nextHref)

		if err != nil {
			log.Fatal("Bad Url", err)
			break
		}

		time.Sleep(1 * time.Second)
	}

	return allProducts, nil
}

func (c *Crawl) fetchAndParse(url string) ([]model.Product, string, error) {
	htmlBody, err := Fetch(url)

	if err != nil {
		return nil, "", err
	}

	return parser.ParseCatalogue(htmlBody, c.BaseUrl)
}

func Fetch(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	return string(body), nil
}
