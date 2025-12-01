package crawler

import (
	"io"
	"net/http"
)

func Fetch(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	return string(body), nil
}
