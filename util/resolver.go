package util

import "net/url"

func Resolve(baseUrl, href string) (string, error) {
	base, err := url.Parse(baseUrl)

	if err != nil {
		return "", err
	}

	rel, err := url.Parse(href)

	if err != nil {
		return "", err
	}

	return base.ResolveReference(rel).String(), nil

}
