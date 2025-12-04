package parser

import (
	"strings"

	"github.com/oloomoses/go-craper/model"
	"github.com/oloomoses/go-craper/util"
	"golang.org/x/net/html"
)

func ParseCatalogue(htmlBody string, baseUrl string) ([]model.Product, string, error) {
	var products []model.Product

	doc, err := html.Parse(strings.NewReader(htmlBody))

	if err != nil {
		return nil, "", err
	}

	var nextUrl string

	var walk func(*html.Node)

	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "article" && hasClass(node, "product_pod") {
			p := extractProduct(node, baseUrl)

			if p != nil {
				products = append(products, *p)
			}
		}

		if node.Type == html.ElementNode && node.Data == "li" && hasClass(node, "next") {
			if a := findChildByTag(node, "a"); a != nil {
				nextUrl = strings.TrimSpace(a.Attr[0].Val)
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)

	return products, nextUrl, nil
}

func hasClass(node *html.Node, cls string) bool {
	for _, a := range node.Attr {
		if a.Key == "class" && a.Val == cls {
			return true
		}
	}
	return false
}

func findChildByTag(n *html.Node, tag string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
	}
	return nil
}

func extractProduct(article *html.Node, baseUrl string) *model.Product {
	titleLink := findChildByTag(article, "h3")

	a := findChildByTag(titleLink, "a")

	title := a.Attr[1].Val

	href := a.Attr[0].Val

	price := getPriceNode(article)

	absUrl, _ := util.Resolve(baseUrl, href)

	return &model.Product{
		Title: title,
		Price: price,
		Url:   absUrl,
	}
}

func getPriceNode(n *html.Node) string {

	var price string

	var f func(*html.Node)

	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "p" && hasClass(node, "price_color") {
			price = strings.TrimSpace(node.FirstChild.Data)
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(n)

	return price
}
