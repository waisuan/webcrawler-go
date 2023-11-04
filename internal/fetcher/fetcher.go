package fetcher

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type IFetcher interface {
	Fetch(targetUrl string) ([]string, error)
}

type Fetcher struct{}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}

func (f *Fetcher) Fetch(rawTargetUrl string) ([]string, error) {
	targetUrl, err := url.Parse(rawTargetUrl)
	if err != nil {
		return nil, err
	}

	content, err := f.getHtmlContent(rawTargetUrl)
	if err != nil {
		return nil, err
	}

	return f.parseAllUrls(content, targetUrl)
}

func (f *Fetcher) getHtmlContent(u string) (string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", err
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return html.UnescapeString(string(content)), err
}

func (f *Fetcher) parseAllUrls(htmlContent string, targetUrl *url.URL) ([]string, error) {
	targetHostname := strings.TrimPrefix(targetUrl.Hostname(), "www.")

	// Search based on the anchor tags in the HTML body.
	re := regexp.MustCompile("<a.*?href=\"(.*?)\"")
	matches := re.FindAllStringSubmatch(htmlContent, -1)

	foundUrls := make(map[string]bool)
	for _, m := range matches {
		foundUrl, err := url.Parse(m[1])
		if err != nil {
			log.Printf("skipping - unable to parse %s\n: %v", m[1], err)
			continue
		}

		// Cater to absolute and relative URLs.
		if foundUrl.IsAbs() {
			// Must match domain of the starting URL.
			foundUrlHostname := strings.TrimPrefix(foundUrl.Hostname(), "www.")
			if foundUrlHostname != targetHostname {
				continue
			}

			foundUrls[foundUrl.String()] = true
		} else {
			foundUrls[targetUrl.Scheme+"://"+targetUrl.Host+foundUrl.String()] = true
		}
	}

	// Dedup found URLs.
	uniqueUrls := make([]string, 0)
	for u := range foundUrls {
		uniqueUrls = append(uniqueUrls, u)
	}

	return uniqueUrls, nil
}
