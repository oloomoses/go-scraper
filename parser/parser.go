package parser

import (
	"strings"

	"golang.org/x/net/html"
)

func ExtractAnkor(htmlBody string) []string {
	doc, err := html.Parse(strings.NewReader(htmlBody))

	if err != nil {
		return []string{}
	}
	var ankor []string
	var visitor func(*html.Node)

	visitor = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, a := range node.Attr {
				if a.Key == "href" {
					ankor = append(ankor, a.Val)
				}
			}
			// return
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			visitor(c)
		}
	}

	visitor(doc)

	return ankor
}
