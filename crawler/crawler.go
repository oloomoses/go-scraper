package crawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/oloomoses/go-craper/model"
	"github.com/oloomoses/go-craper/parser"
	"github.com/oloomoses/go-craper/polite"
	"github.com/oloomoses/go-craper/util"
)

type Crawl struct {
	BaseUrl string
	Visited map[string]bool
	Robots  *polite.Robots
}

func New(baseUrl string) *Crawl {
	return &Crawl{
		BaseUrl: baseUrl,
		Visited: make(map[string]bool),
		Robots:  polite.NewRobots(baseUrl),
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

		nextFullUrl, err := util.Resolve(currentUrl, nextHref)

		if err != nil {
			log.Fatal("Bad Url", err)
			break
		}

		parsedUrl, _ := url.Parse(nextFullUrl)

		if !c.Robots.CanFetch(parsedUrl.Path) {
			log.Printf("Blocked by robots.txt: %s", parsedUrl.Path)
			break
		}

		currentUrl = nextFullUrl

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
