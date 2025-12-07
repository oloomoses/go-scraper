package crawler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sync"
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

func (c *Crawl) CrawlConcurrently(ctx context.Context, baseUrl string, maxPages int, workers int) ([]model.Product, error) {
	var allProducts []model.Product
	urls, err := c.DiscoverPages(baseUrl, maxPages)
	wg := &sync.WaitGroup{}

	if err != nil {
		return nil, err
	}

	jobs := make(chan string, len(urls))

	results := make(chan []model.Product, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			for u := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}
				products, _, err := c.fetchAndParse(u)

				if err != nil {
					log.Printf("Worker %d failed %s: %v \n", id, u, err)
					continue
				}
				results <- products
			}
		}(i + 1)
	}

	go func() {
		defer close(jobs)
		for _, u := range urls {
			select {
			case jobs <- u:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer close(results)
		for prods := range results {
			allProducts = append(allProducts, prods...)
		}

	}()

	wg.Wait()

	return allProducts, nil
}

func (c *Crawl) fetchAndParse(url string) ([]model.Product, string, error) {
	htmlBody, err := Fetch(url)

	if err != nil {
		return nil, "", err
	}

	return parser.ParseCatalogue(htmlBody, c.BaseUrl)
}

func (c *Crawl) DiscoverPages(baseUrl string, maxPages int) ([]string, error) {
	var urls []string

	urls = append(urls, baseUrl)
	pageNumber := 0

	newBaseUrl, err := url.Parse(baseUrl)

	if err != nil {
		return nil, err
	}

	if newBaseUrl.Path == "" {
		pageNumber = 1
	} else {

		re := regexp.MustCompile(`page-(\d+)`)
		match := re.FindStringSubmatch(newBaseUrl.Path)

		fmt.Sscanf(match[1], "%d", &pageNumber)
	}

	for i := 1; i <= maxPages; i++ {
		pageNumber++
		u := fmt.Sprintf("/catalogue/page-%d.html", pageNumber)

		resolvedUrl, _ := util.Resolve(baseUrl, u)

		urls = append(urls, resolvedUrl)
	}

	return urls, nil
}

func Fetch(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 50,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false,
			ForceAttemptHTTP2:   true,
		},
	}
	res, err := client.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	return string(body), nil
}
