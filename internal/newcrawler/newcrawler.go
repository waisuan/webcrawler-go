package newcrawler

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
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

var wg sync.WaitGroup

func (c *Crawler) Run(url string, depth int) {
	consume := func(id int, targetUrlCh <-chan string, pendingUrlsCh chan<- []string) {
		defer wg.Done()

		for {
			url := <-targetUrlCh

			o := c.markAsVisited(url)
			if !o {
				continue
			}

			log.Printf("visited: %s\n", url)

			urls, err := c.fetcher.Fetch(url)
			if err != nil {
				log.Printf("skipping - unable to crawl %s - %v\n", url, err)
				return
			}

			if len(urls) == 0 {
				continue
			}

			var visitingUrls string
			if len(urls) > c.cfg.MaxLoggedUrls {
				visitingUrls = fmt.Sprintf("%d links", len(urls))
			} else {
				visitingUrls = strings.Join(urls, ", ")
			}
			log.Printf("\twill try visiting: %s\n", visitingUrls)

			pendingUrlsCh <- urls
		}
	}

	produce := func(pendingUrlsCh chan []string, targetUrlCh chan<- string) {
		attempts := 0
		for {
			select {
			case pendingUrls := <-pendingUrlsCh:
				for _, u := range pendingUrls {
					targetUrlCh <- u
					attempts = 0
				}
			case <-time.After(1 * time.Second):
				log.Println("searching for more links to process...")
			}

			if attempts > 10 {
				break
			}

			attempts += 1
		}

		close(pendingUrlsCh)
		close(targetUrlCh)
	}

	targetUrlCh := make(chan string)

	pendingUrlsCh := make(chan []string, 1000)

	for i := 0; i < c.cfg.MaxCrawlConcurrencyLevel; i++ {
		wg.Add(1)
		go consume(i, targetUrlCh, pendingUrlsCh)
	}

	go produce(pendingUrlsCh, targetUrlCh)

	pendingUrlsCh <- []string{url}

	wg.Wait()
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
