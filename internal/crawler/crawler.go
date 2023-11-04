package crawler

import (
	"fmt"
	"log"
	"sync"
	"webcrawler-go/internal/fetcher"
)

type Crawler struct {
	fetcher fetcher.IFetcher
	visited map[string]bool
	lock    sync.Mutex
}

func NewCrawler(fetcher fetcher.IFetcher) *Crawler {
	return &Crawler{
		fetcher: fetcher,
		visited: make(map[string]bool),
	}
}

func (c *Crawler) Run(url string) {
	o := c.markAsVisited(url)
	if !o {
		return
	}

	urls, err := c.fetcher.Fetch(url)
	if err != nil {
		log.Printf("unable to fetch from %s - %v\n", url, err)
		return
	}

	fmt.Printf("found: %s\n", url)

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			c.Run(u)
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
