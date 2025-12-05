package polite

import (
	"bufio"
	"io"
	"net/http"
	"strings"
	"time"
)

type Robots struct {
	disallowed []string
	canFetch   func(path string) bool
}

func NewRobots(baseUrl string) *Robots {
	r := &Robots{
		disallowed: make([]string, 0),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Get(baseUrl + "/robots.txt")

	if err != nil || res.StatusCode != 200 {
		if res != nil {
			res.Body.Close()
		}
		r.canFetch = func(_ string) bool { return true }

		return r
	}

	defer res.Body.Close()

	r.parse(res.Body)

	return r
}

func (r *Robots) parse(body io.Reader) {
	scanner := bufio.NewScanner(body)

	active := false

	r.canFetch = func(path string) bool {
		for _, dis := range r.disallowed {
			if strings.HasPrefix(path, dis) {
				return false
			}
		}
		return true
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "User-agent:") {
			active = strings.Contains(line, "*") || strings.Contains(line, "User-Agent: *")
			continue
		}

		if active && strings.HasPrefix(line, "Disallow:") {
			parts := strings.SplitN(line, ":", 2)

			if len(parts) == 2 {
				path := strings.TrimSpace(parts[1])

				if path != "" && path != "/" {
					r.disallowed = append(r.disallowed, path)
				}
			}
		}
	}

}

func (r *Robots) CanFetch(path string) bool {
	return r.canFetch(path)
}
