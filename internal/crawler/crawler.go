package crawler

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"webcrawler-go/internal/fetcher"
)

const MaxVisitingUrlsToPrint = 20

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
	if len(urls) > MaxVisitingUrlsToPrint {
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
