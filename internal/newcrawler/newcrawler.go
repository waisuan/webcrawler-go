package newcrawler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"webcrawler-go/internal/dependencies"
	"webcrawler-go/internal/fetcher"
)

var (
	wg   sync.WaitGroup
	work atomic.Int64
)

type Crawler struct {
	cfg     *dependencies.Config
	fetcher fetcher.IFetcher
	visited map[string]bool
	lock    sync.Mutex
}

type CrawlUrlJob struct {
	url   string
	depth int
}

func NewCrawler(cfg *dependencies.Config, fetcher fetcher.IFetcher) *Crawler {
	return &Crawler{
		cfg:     cfg,
		fetcher: fetcher,
		visited: make(map[string]bool),
	}
}

func (c *Crawler) Run(url string, depth int) {
	worker := func(ctx context.Context, targetUrlCh <-chan string, pendingUrlsCh chan<- []string) {
		defer wg.Done()

	loop:
		for {
			select {
			case url := <-targetUrlCh:
				work.Add(1)

				o := c.markAsVisited(url)
				if !o {
					work.Add(-1)
					continue
				}

				log.Printf("visited: %s\n", url)

				urls, err := c.fetcher.Fetch(url)
				if err != nil {
					log.Printf("skipping - unable to crawl %s - %v\n", url, err)
					work.Add(-1)
					continue
				}

				if len(urls) == 0 {
					work.Add(-1)
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

				work.Add(-1)
			case <-ctx.Done():
				break loop
			}
		}
	}

	crawl := func(terminator context.CancelFunc, targetUrlCh chan<- string, pendingUrlsCh chan []string) {
	loop:
		for {
			select {
			case pendingUrls := <-pendingUrlsCh:
				for _, u := range pendingUrls {
					targetUrlCh <- u
				}
				continue
			case <-time.After(5 * time.Second): // helps to terminate all workers when there's nothing left to process.
				log.Printf("searching for more links to process...(%d)\n", work.Load())
				if work.Load() <= 0 {
					break loop
				}
			}
		}

		terminator()
	}

	targetUrlCh := make(chan string)
	defer close(targetUrlCh)

	pendingUrlsCh := make(chan []string, 10_000) // Buffered channel to limit the no. of pending unprocessed links at a time.
	defer close(pendingUrlsCh)

	ctx, terminator := context.WithCancel(context.Background())

	for i := 0; i < c.cfg.MaxCrawlConcurrencyLevel; i++ {
		wg.Add(1)
		go worker(ctx, targetUrlCh, pendingUrlsCh)
	}

	go crawl(terminator, targetUrlCh, pendingUrlsCh)

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
