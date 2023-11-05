package crawler

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"webcrawler-go/internal/dependencies"
	"webcrawler-go/internal/fetcher"
)

type Crawler struct {
	cfg     *dependencies.Config
	fetcher fetcher.IFetcher
	visited map[string]bool
	lock    sync.Mutex
}

func NewCrawler(cfg *dependencies.Config, fetcher fetcher.IFetcher) *Crawler {
	return &Crawler{
		cfg:     cfg,
		fetcher: fetcher,
		visited: make(map[string]bool),
	}
}

func (c *Crawler) Run(url string, depth int) {
	o := c.markAsVisited(url)
	if !o || (c.cfg.MaxCrawlDepth > 0 && depth >= c.cfg.MaxCrawlDepth) {
		return
	}

	log.Printf("visited: %s\n", url)

	urls, err := c.fetcher.Fetch(url)
	if err != nil {
		log.Printf("skipping - unable to crawl %s - %v\n", url, err)
		return
	}

	if len(urls) == 0 {
		return
	}

	var visitingUrls string
	if len(urls) > c.cfg.MaxLoggedUrls {
		visitingUrls = fmt.Sprintf("%d links", len(urls))
	} else {
		visitingUrls = strings.Join(urls, ", ")
	}
	log.Printf("\twill try visiting: %s\n", visitingUrls)

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			c.Run(u, depth+1)
		}(u)
	}
	wg.Wait()

	return
}

func (c *Crawler) markAsVisited(url string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.visited[url]
	if ok {
		return false
	}
	c.visited[url] = true

	return true
}
